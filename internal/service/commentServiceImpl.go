package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/IRONICBo/QiYin_BE/internal/conn/db"
	"github.com/IRONICBo/QiYin_BE/internal/dal/dao"
	"github.com/IRONICBo/QiYin_BE/internal/middleware/rabbitmq"
	"github.com/IRONICBo/QiYin_BE/internal/utils"
	"github.com/gin-gonic/gin"
)

type CommentServiceImpl struct {
	UserService
}

func NewCommentService(c *gin.Context) *CommentServiceImpl {
	return &CommentServiceImpl{
		UserService: &UserServiceImpl{},
	}
}

// CountFromVideoId  获取视频的评论数.
func (c CommentServiceImpl) CountFromVideoId(videoId int64) (int64, error) {
	key := fmt.Sprintf("%s:%d", utils.Comment, videoId)
	ctx := context.Background()
	cnt, err := db.GetRedis().SCard(ctx, key).Result()
	if err != nil { // 若查询缓存出错，则打印log
		// return 0, err
		log.Println("count from redis error:", err)
	}
	if cnt != 0 {
		return cnt - 1, nil
	}

	// 缓存中查不到则去数据库查  并存储在数据库中
	cntDao, err1 := dao.Count(videoId)
	if err1 != nil {
		log.Println("comment count dao err:", err1)
		return 0, nil
	}
	// 将评论id切片存入redis-第一次存储 V-C set 值：
	go func() {
		// 查询评论id list
		cList, _ := dao.CommentIdList(videoId)
		// 先在redis中存储一个-1值，防止脏读
		_, _err := db.GetRedis().SAdd(ctx, key, string(rune(utils.DefaultRedisValue))).Result()
		if _err != nil { // 若存储redis失败，则直接返回
			log.Println("redis save one vId - cId 0 failed")
			return
		}
		// 设置key值过期时间
		_, err := db.GetRedis().Expire(ctx, key,
			time.Duration(utils.OneMonth)*time.Second).Result()
		if err != nil {
			log.Println("redis save one vId - cId expire failed")
		}
		// 评论id循环存入redis
		for _, commentId := range cList {
			insertRedisVideoCommentId(strconv.Itoa(int(videoId)), commentId)
		}
		log.Println("count comment save ids in redis")
	}()
	// 返回结果
	return cntDao, nil
}

func (c CommentServiceImpl) CommentAdd(comment dao.Comment) (dao.CommentData, error) {
	// 1.评论信息存储：
	commentId, err := dao.InsertComment(comment)
	if err != nil {
		return dao.CommentData{}, err
	}
	// 2.查询用户信息
	userData, err2 := c.GetTableUserById(comment.UserId)
	if err2 != nil {
		return dao.CommentData{}, err2
	}
	// 3.拼接
	commentData := dao.CommentData{
		Id:         commentId,
		UserInfo:   userData,
		Content:    comment.CommentText,
		CreateDate: comment.CreateDate.Format(utils.DateTime),
	}
	// 将此发表的评论id存入redis
	go func() {
		insertRedisVideoCommentId(strconv.Itoa(int(comment.VideoId)), strconv.Itoa(int(comment.Id)))
		log.Println("send comment save in redis")
	}()
	// 返回结果
	return commentData, nil
}

func (c CommentServiceImpl) CommentDelete(userId string, commentId int64) error {
	// 首先判断一下有没有删除权限 去数据库找一下该评论的 userid是不是这个
	err := dao.GetUserIdByCommentId(commentId, userId)
	if err != nil {
		log.Println(err)
	}

	key := fmt.Sprintf("%s:%d", utils.CommentCV, commentId)
	ctx := context.Background()
	// 1.先查询redis，若有则删除，返回客户端-再go协程删除数据库；无则在数据库中删除，返回客户端。
	n, err := db.GetRedis().Exists(ctx, key).Result()
	if err != nil {
		log.Println(err)
	}
	if n > 0 { // 在缓存中有此值，则找出来删除，然后返回
		vid, err1 := db.GetRedis().Get(ctx, key).Result()
		if err1 != nil { // 没找到，返回err
			log.Println("redis find CV err:", err1)
		}
		// 删除，	删除了CV  下一个删除VC
		del1, err2 := db.GetRedis().Del(ctx, key).Result()
		if err2 != nil {
			log.Println(err2)
		}

		del2, err3 := db.GetRedis().SRem(ctx, fmt.Sprintf("%s:%s", utils.Comment, vid), strconv.FormatInt(commentId, 10)).Result()
		if err3 != nil {
			log.Println(err3)
		}
		log.Println("del comment in Redis success:", del1, del2) // del1、del2代表删除了几条数据
		// 使用mq进行数据库中评论的删除-评论状态更新
		// 评论id传入消息队列
		rabbitmq.RmqCommentDel.Publish(strconv.FormatInt(commentId, 10))
		return nil
	}
	// 不在内存中，则直接走数据库删除
	return dao.DeleteComment(commentId)
}

// GetList.
func (c CommentServiceImpl) GetList(videoId string) ([]dao.CommentData, error) {
	// 1.先查询评论列表信息
	commentList, err := dao.GetCommentList(videoId)
	if err != nil {
		log.Println("CommentService-GetList: return err: " + err.Error()) // 函数返回提示错误信息
		return nil, err
	}
	// 当前有0条评论
	if commentList == nil {
		return nil, nil
	}

	// 提前定义好切片长度
	commentInfoList := make([]dao.CommentData, len(commentList))

	wg := &sync.WaitGroup{}
	wg.Add(len(commentList))
	idx := 0
	for _, comment := range commentList {
		// 2.调用方法组装评论信息，再append
		var commentData dao.CommentData
		// 将评论信息进行组装，添加想要的信息,插入从数据库中查到的数据
		go func(comment dao.Comment) {
			c.oneComment(&commentData, &comment)
			// 3.组装list
			// commentInfoList = append(commentInfoList, commentData)
			commentInfoList[idx] = commentData
			idx = idx + 1
			wg.Done()
		}(comment)
	}
	wg.Wait()
	////评论排序-按照主键排序
	//sort.Sort(CommentSlice(commentInfoList))
	key := fmt.Sprintf("%s:%s", utils.Comment, videoId)
	ctx := context.Background()
	// 协程查询redis中是否有此记录，无则将评论id切片存入redis
	go func() {
		// 1.先在缓存中查此视频是否已有评论列表
		cnt, err1 := db.GetRedis().SCard(ctx, key).Result()
		if err1 != nil { // 若查询缓存出错，则打印log
			// return 0, err
			log.Println("count from redis error:", err)
		}
		// 2.缓存中查到了数量大于0，则说明数据正常，不用更新缓存
		if cnt > 1 {
			return
		}
		// 3.缓存中数据不正确，更新缓存：
		// 先在redis中存储一个-1 值，防止脏读
		_, _err := db.GetRedis().SAdd(ctx, key, string(rune(utils.DefaultRedisValue))).Result()
		if _err != nil { // 若存储redis失败，则直接返回
			log.Println("redis save one vId - cId 0 failed")
			return
		}
		// 设置key值过期时间
		_, err2 := db.GetRedis().Expire(ctx, key,
			time.Duration(utils.OneMonth)*time.Second).Result()
		if err2 != nil {
			log.Println("redis save one vId - cId expire failed")
		}
		// 将评论id循环存入redis
		for _, _comment := range commentInfoList {
			insertRedisVideoCommentId(videoId, strconv.Itoa(int(_comment.Id)))
		}
		log.Println("comment list save ids in redis")
	}()

	log.Println("CommentService-GetList: return list success") // 函数执行成功，返回正确信息
	return commentInfoList, nil
}

// 在redis中存储video_id对应的comment_id 、 comment_id对应的video_id.
func insertRedisVideoCommentId(videoId string, commentId string) {
	key := fmt.Sprintf("%s:%s", utils.Comment, videoId)
	ctx := context.Background()
	// 在redis-RdbVCid中存储video_id对应的comment_id
	_, err := db.GetRedis().SAdd(ctx, key, commentId).Result()
	if err != nil { // 若存储redis失败-有err，则直接删除key
		log.Println("redis save send: vId - cId failed, key deleted")
		db.GetRedis().Del(ctx, videoId)
		return
	}
	keyReverse := fmt.Sprintf("%s:%s", utils.CommentCV, commentId)
	// 在redis-RdbCVid中存储comment_id对应的video_id
	_, err = db.GetRedis().Set(ctx, keyReverse, videoId, 0).Result()
	if err != nil {
		log.Println("redis save one cId - vId failed")
	}
}

// 此函数用于给一个评论赋值：评论信息+用户信息 填充.
func (c CommentServiceImpl) oneComment(comment *dao.CommentData, com *dao.Comment) {
	var wg sync.WaitGroup
	wg.Add(1)
	// 根据评论用户id和当前用户id，查询评论用户信息
	var err error
	comment.Id = com.Id
	comment.Content = com.CommentText
	comment.CreateDate = com.CreateDate.Format(utils.DateTime)
	comment.UserInfo, err = c.GetTableUserById(com.UserId)
	if err != nil {
		log.Println("CommentService-GetList: GetUserByIdWithCurId return err: " + err.Error()) // 函数返回提示错误信息
	}
	wg.Done()
	wg.Wait()
}

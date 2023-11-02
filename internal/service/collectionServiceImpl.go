package service

import (
	"context"
	"fmt"
	"github.com/IRONICBo/QiYin_BE/internal/conn/db"
	"github.com/IRONICBo/QiYin_BE/internal/dal/dao"
	"github.com/IRONICBo/QiYin_BE/internal/middleware/rabbitmq"
	"github.com/IRONICBo/QiYin_BE/internal/utils"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type CollectionServiceImpl struct {
	VideoService
	UserService
}

// NewCollectionService return new service with gin context.
func NewCollectionService(c *gin.Context) *CollectionServiceImpl {
	return &CollectionServiceImpl{
		UserService: &UserServiceImpl{
			Service: Service{
				ctx: c,
			},
		},
		VideoService: &VideoServiceImpl{},
	}
}

// IsCollection 根据collection:userId,videoId查询点赞状态;
func (collection *CollectionServiceImpl) IsCollection(videoId string, userId string) (bool, error) {
	key := fmt.Sprintf("%s:%s", utils.Collection, userId)
	ctx := context.Background()
	if n, err := db.GetRedis().Exists(ctx, key).Result(); n > 1 {
		if err != nil {
			log.Printf("this user has no collection：%v", err)
			return false, err
		}
		exist, err1 := db.GetRedis().SIsMember(ctx, key, videoId).Result()
		if err1 != nil {
			log.Printf("This user didn't collection this video：%v", err1)
			return false, err1
		}
		log.Printf("IsCollection query success")
		return exist, nil
	} else {
		//反过来查一下
		reverseKey := fmt.Sprintf("%s:%s", utils.Collection, videoId)
		if n, err := db.GetRedis().Exists(ctx, reverseKey).Result(); n > 1 {
			if err != nil {
				log.Printf("this video has no collection：%v", err)
				return false, err
			}
			exist, err1 := db.GetRedis().SIsMember(ctx, reverseKey, userId).Result()
			if err1 != nil {
				log.Printf("This video didn't collection this user：%v", err1)
				return false, err1
			}
			log.Printf("IsCollection query success")
			return exist, nil
		} else {
			//两种方式都没有查到  所以都需要进行同步
			_, ok, err := collection.collectionToRedis(ctx, key, userId, false)
			if err != nil || !ok {
				return false, err
			}
			_, ok, err = collection.collectionToRedis(ctx, fmt.Sprintf("%s:%s", utils.Collection, videoId), videoId, true)
			if err != nil || !ok {
				return false, err
			}
			//查询Redis CollectionUserId,key：strUserId中是否存在value:videoId,存在返回true,不存在返回false
			exist, err2 := db.GetRedis().SIsMember(ctx, key, videoId).Result()
			//如果有问题，说明操作redis失败,返回默认false,返回错误信息
			if err2 != nil {
				return false, err2
			}
			return exist, nil
		}
	}
}

// 同步点赞状态mysql到redis
func (collection *CollectionServiceImpl) collectionToRedis(ctx context.Context, key string, userId string, isReverse bool) ([]string, bool, error) {
	//如果不存在，则维护Redis CollectionUserId 新建key:key,设置过期时间，加入DefaultRedisValue，
	//key:key，加入value:DefaultRedisValue,过期才会删，防止删最后一个数据的时候数据库还没更新完出现脏读，或者数据库操作失败造成的脏读
	//通过userId查询collections表,返回所有点赞videoId，加入key:strUserId集合中,
	//再加入当前videoId,再更新collections表此条数据
	if _, err := db.GetRedis().SAdd(ctx, key, string(rune(utils.DefaultRedisValue))).Result(); err != nil {
		log.Print(err)
		db.GetRedis().Del(ctx, key)
		return []string{}, false, err
	}
	_, err := db.GetRedis().Expire(ctx, key, time.Duration(utils.OneMonth)*time.Second).Result()
	if err != nil {
		log.Printf("set expire failed")
		db.GetRedis().Del(ctx, key)
		return []string{}, false, err
	}
	var IdList []string
	// video 作为id
	if isReverse {
		videoId, _ := strconv.ParseInt(userId, 10, 64)
		IdList, err = dao.GetCollectionUserIdList(videoId)
	} else {
		IdList, err = dao.GetCollectionVideoIdList(userId)
	}

	//videoIdList, err1 := dao.GetCollectionVideoIdList(userId)
	//如果有问题，说明查询失败，返回错误信息："get collectionVideoIdList failed"
	if err != nil {
		return []string{}, false, err
	}

	//维护Redis CollectionUserId(key:key)，遍历videoIdList加入
	for _, collectionId := range IdList {
		if _, err1 := db.GetRedis().SAdd(ctx, key, collectionId).Result(); err1 != nil {
			db.GetRedis().Del(ctx, key)
			log.Println("async failed")
			return []string{}, false, err1
		}
	}
	return IdList, true, nil
}

// 点赞
func (collection *CollectionServiceImpl) collectionDo(key string, videoId string, userId string, sb strings.Builder, isReverse bool) error {
	ctx := context.Background()
	//查询Redis CollectionUserId(key:key)是否已经加载过此信息
	if n, err := db.GetRedis().Exists(ctx, key).Result(); n > 0 {
		//如果有问题，说明查询redis失败,返回错误信息
		if err != nil {
			log.Printf("CollectionAction query failed：%v", err)
			return err
		} //如果加载过此信息key:key，则加入value:videoId
		if _, err1 := db.GetRedis().SAdd(ctx, key, videoId).Result(); err1 != nil {
			log.Print(err1)
			return err1
		} else {
			rabbitmq.RmqCollectionAdd.Publish(sb.String())
		}
	} else {
		_, _, err := collection.collectionToRedis(ctx, key, userId, isReverse)
		if err != nil {
			return err
		}
		if _, err2 := db.GetRedis().SAdd(ctx, key, videoId).Result(); err2 != nil {
			return err2
		} else {

			rabbitmq.RmqCollectionAdd.Publish(sb.String())
		}
	}

	return nil
}

func (collection *CollectionServiceImpl) collectionCancel(key string, videoId string, userId string, sb strings.Builder, isReverse bool) error {
	ctx := context.Background()
	//查询Redis CollectionUserId(key:key)是否已经加载过此信息
	if n, err := db.GetRedis().Exists(ctx, key).Result(); n > 0 {
		//如果有问题，说明查询redis失败,返回错误信息
		if err != nil {
			return err
		} //防止出现redis数据不一致情况，当redis删除操作成功，才执行数据库更新操作
		if _, err1 := db.GetRedis().SRem(ctx, key, videoId).Result(); err1 != nil {
			return err1
		} else {
			//后续数据库的操作，可以在mq里设置若执行数据库更新操作失败，重新消费该信息

			rabbitmq.RmqCollectionDel.Publish(sb.String())
		}
	} else {
		_, _, err := collection.collectionToRedis(ctx, key, userId, isReverse)
		if err != nil {
			return err
		}

		if _, err2 := db.GetRedis().SRem(ctx, key, videoId).Result(); err2 != nil {
			return err2
		} else {

			rabbitmq.RmqCollectionDel.Publish(sb.String())
		}
	}
	return nil
}

// CollectionAction 根据userId，videoId,actionType对视频进行点赞或者取消赞操作;
// step1: 维护Redis CollectionUserId(key:key),添加或者删除value:videoId,CollectionVideoId(key:strVideoId),添加或者删除value:userId;
// step2：更新数据库collections表;
func (collection *CollectionServiceImpl) CollectionAction(userId string, videoId string, actionType int32) error {
	key := fmt.Sprintf("%s:%s", utils.Collection, userId)
	//将要操作数据库collections表的信息打入消息队列RmqCollectionAdd或者RmqCollectionDel
	//拼接打入信息
	sb := strings.Builder{}
	sb.WriteString(userId)
	sb.WriteString(" ")
	sb.WriteString(videoId)

	//执行点赞操作维护
	if actionType == utils.IsCollection {
		err := collection.collectionDo(key, videoId, userId, sb, false)
		if err != nil {
			return err
		}
		err = collection.collectionDo(fmt.Sprintf("%s:%s", utils.Collection, videoId), userId, videoId, sb, true)
		if err != nil {
			return err
		}
	} else { //执行取消赞操作维护
		err := collection.collectionCancel(key, videoId, userId, sb, false)
		if err != nil {
			return err
		}
		err = collection.collectionCancel(fmt.Sprintf("%s:%s", utils.Collection, videoId), userId, videoId, sb, true)
		if err != nil {
			return err
		}
	}
	return nil
}

// TotalCollection 根据userId获取这个用户总共被点赞数量
// 首先需要找到用户的视频  并统计每个视频的被点赞数
func (collection *CollectionServiceImpl) TotalCollection(userId string) (int64, error) {
	//根据userId获取这个用户的发布视频列表信息
	videoIdList, err := collection.GetVideoIdList(userId)
	if err != nil {
		return 0, err
	}
	var sum int64
	//提前开辟空间,存取每个视频的点赞数
	videoCollectionCountList := new([]int64)
	//采用协程并发将对应videoId的点赞数添加到集合中去
	i := len(videoIdList)
	var wg sync.WaitGroup
	wg.Add(i)
	for j := 0; j < i; j++ {
		go collection.addVideoCollectionCount(videoIdList[j], videoCollectionCountList, &wg)
	}
	wg.Wait()
	//遍历累加，求总被点赞数
	for _, count := range *videoCollectionCountList {
		sum += count
	}
	return sum, nil
}

// CollectionCount 根据videoId获取对应点赞数量;
// step1：查询Redis CollectionVideoId(key:strVideoId)是否已经加载过此信息，通过set集合中userId个数，获取点赞数量;
// step2：CollectionVideoId中都没有对应key，维护CollectionVideoId对应key，再通过set集合中userId个数，获取点赞数量;
func (collection *CollectionServiceImpl) CollectionCount(videoId string, isReverse bool) (int64, error) {
	//将int64 videoId转换为 string strVideoId
	key := fmt.Sprintf("%s:%s", utils.Collection, videoId)
	ctx := context.Background()

	//step1 如果key:strVideoId存在 则计算集合中userId个数   set中会有默认值  所以必须要大于1
	if n, err := db.GetRedis().Exists(ctx, key).Result(); n > 1 {
		//如果有问题，说明查询redis失败,返回默认false,返回错误信息
		if err != nil {
			log.Printf("CollectionCount query failed：%v", err)
			return 0, err
		}
		//获取集合中userId个数
		count, err1 := db.GetRedis().SCard(ctx, videoId).Result()
		//如果有问题，说明操作redis失败,返回默认0,返回错误信息
		if err1 != nil {
			return 0, err1
		}
		return count - 1, nil //去掉默认的
	} else {
		_, _, err1 := collection.collectionToRedis(ctx, key, videoId, isReverse)
		if err1 != nil {
			return 0, err1
		}
		//再通过set集合中userId个数,获取点赞数量
		count, err2 := db.GetRedis().SCard(ctx, key).Result()
		//fmt.Println(key)
		//fmt.Println(db.GetRedis().SMembers(ctx, key))
		//fmt.Println(count)
		//如果有问题，说明操作redis失败,返回默认0,返回错误信息
		if err2 != nil {
			log.Printf("CollectionCount query count failed%v", err2)
			return 0, err2
		}
		return count - 1, nil
	}
}

// CollectionVideoCount 根据userId获取这个用户点赞视频数量
func (collection *CollectionServiceImpl) CollectionVideoCount(userId string) (int64, error) {
	count, err := collection.CollectionCount(userId, false)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// addVideoCollectionCount 根据videoId，将该视频点赞数加入对应提前开辟好的空间内
func (collection *CollectionServiceImpl) addVideoCollectionCount(videoId int64, videoCollectionCountList *[]int64, wg *sync.WaitGroup) {
	defer wg.Done()
	//调用CollectionCount：根据videoId,获取点赞数
	count, err := collection.CollectionCount(strconv.FormatInt(videoId, 10), true)
	if err != nil {
		//如果有错误，输出错误信息，并不加入该视频点赞数
		log.Printf(err.Error())
		return
	}
	*videoCollectionCountList = append(*videoCollectionCountList, count)
}

// CollectionList 根据userId，curId(当前用户Id),返回userId的点赞列表;
// step1：查询Redis CollectionUserId(key:key)是否已经加载过此信息，获取集合中全部videoId，并添加到点赞列表集合中;
// step2：CollectionUserId中都没有对应key，维护CollectionUserId对应key，同时添加到点赞列表集合中;
func (collection *CollectionServiceImpl) CollectionList(userId string, curId string) ([]dao.ResVideo, error) {
	ctx := context.Background()
	//将int64 userId转换为 string key
	key := fmt.Sprintf("%s:%s", utils.Collection, userId)
	//step1:查询Redis CollectionUserId,如果key：strUserId存在,则获取集合中全部videoId
	if n, err := db.GetRedis().Exists(ctx, key).Result(); n > 1 {
		//如果有问题，说明查询redis失败,返回默认nil,返回错误信息
		if err != nil {
			log.Printf("CollectionList query failed：%v", err)
			return nil, err
		}
		//获取集合中全部videoId
		videoIdList, err1 := db.GetRedis().SMembers(ctx, key).Result()
		//如果有问题，说明查询redis失败,返回默认nil,返回错误信息
		if err1 != nil {
			log.Printf("CollectionList get values failed：%v", err1)
			return nil, err1
		}
		//提前开辟点赞列表空间
		collectionVideoList := new([]dao.ResVideo)
		//采用协程并发将Video类型对象添加到集合中去
		i := len(videoIdList) - 1 //去掉DefaultRedisValue
		if i == 0 {
			return *collectionVideoList, nil
		}
		var wg sync.WaitGroup
		wg.Add(i)
		for j := 0; j <= i; j++ {
			//将string videoId转换为 int64 VideoId
			videoId, _ := strconv.ParseInt(videoIdList[j], 10, 64)
			if videoId == utils.DefaultRedisValue {
				continue
			}
			go collection.addCollectionVideoList(videoId, curId, collectionVideoList, &wg)
		}
		wg.Wait()
		return *collectionVideoList, nil
	} else {
		ids, _, err := collection.collectionToRedis(ctx, key, userId, false)
		if err != nil {
			return nil, err
		}
		//提前开辟点赞列表空间
		collectionVideoList := new([]dao.ResVideo)
		i := len(ids)
		if i == 0 {
			return *collectionVideoList, nil
		}
		var wg sync.WaitGroup
		wg.Add(i)
		for j := 0; j < i; j++ {
			videoId, _ := strconv.ParseInt(ids[j], 10, 64)
			go collection.addCollectionVideoList(videoId, curId, collectionVideoList, &wg)
		}
		wg.Wait()
		return *collectionVideoList, nil
	}
}

// addCollectionVideoList 根据videoId,登录用户curId，添加视频对象到点赞列表空间
func (collection *CollectionServiceImpl) addCollectionVideoList(videoId int64, curId string, collectionVideoList *[]dao.ResVideo, wg *sync.WaitGroup) {
	defer wg.Done()
	//调用videoService接口，GetVideo：根据videoId，当前用户id:curId，返回Video类型对象
	video, err := collection.GetVideo(videoId, curId)
	if err != nil {
		log.Println("this collection video is miss")
		return
	}
	//将Video类型对象添加到集合中去
	*collectionVideoList = append(*collectionVideoList, video)
}

func GetCollectionService() CollectionServiceImpl {
	var userService UserServiceImpl
	var videoService VideoServiceImpl
	var favoriteService FavoriteServiceImpl
	var collectionService CollectionServiceImpl
	userService.FavoriteService = &favoriteService
	userService.CollectionService = &collectionService
	favoriteService.VideoService = &videoService
	collectionService.VideoService = &videoService
	videoService.UserService = &userService
	return collectionService
}

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
	"github.com/IRONICBo/QiYin_BE/internal/utils"
	"github.com/gin-gonic/gin"
)

type VideoServiceImpl struct {
	UserService
	FavoriteService
	CollectionService
	CommentService
}

// NewVideoService return new service with gin context.
func NewVideoService(c *gin.Context) *VideoServiceImpl {
	return &VideoServiceImpl{
		FavoriteService:   &FavoriteServiceImpl{},
		CollectionService: &CollectionServiceImpl{},
		UserService:       &UserServiceImpl{},
		CommentService: &CommentServiceImpl{
			UserService: &UserServiceImpl{},
		},
	}
}

// GetVideoIdList
// 通过一个作者id，返回该用户发布的视频id切片数组.
func (videoService *VideoServiceImpl) GetVideoIdList(userId string) ([]int64, error) {
	// 直接调用dao层方法获取id即可
	id, err := dao.GetVideoIdsByAuthorId(userId)
	if err != nil {
		return nil, err
	}
	return id, nil
}

// Search 搜索.
func (videoService *VideoServiceImpl) Search(searchValue string, userId string) ([]dao.ResVideo, error) {
	// 查询到相关的videolist + 相关的用户信息
	videoList, err := dao.GetVideoByTitle(searchValue)
	// 查询失败直接返回
	if err != nil {
		log.Printf("Search query failed：%v", err)
		return videoList, err
	}

	//	将查询词放到redis中以便热榜搜索
	ctx := context.Background()
	key := utils.Search
	if _, err := db.GetRedis().ZIncrBy(ctx, key, 1, searchValue).Result(); err != nil {
		log.Printf("ZincrBt failed：%v", err)
	}

	// 设置过期时间
	_, err = db.GetRedis().Expire(ctx, key, time.Duration(utils.OneDay)*time.Second).Result()
	if err != nil {
		fmt.Println("Error:", err)
	}

	// 得到点赞数和收藏数
	for _, video := range videoList {
		videoService.creatVideo(&video, userId)
	}

	return videoList, nil
}

// GetHots 得到热榜.
func (videoService *VideoServiceImpl) GetHots() ([]string, error) {
	ctx := context.Background()
	res, err := db.GetRedis().ZRevRange(ctx, utils.Search, 0, 9).Result()
	if err != nil {
		return []string{}, err
	}
	return res, nil
}

// GetVideo 根据videoId 和 curUserId 得到视频
func (videoService *VideoServiceImpl) GetVideo(videoId int64, curUserId string) (dao.ResVideo, error) {
	// 查询到相关的videolist + 相关的用户信息
	video, err := dao.GetVideoById(videoId)
	// 查询失败直接返回
	if err != nil {
		log.Printf("GetVideo failed：%v", err)
		return video, err
	}

	videoService.creatVideo(&video, curUserId)
	return video, nil
}

// 将video进行组装，添加想要的信息,插入从数据库中查到的数据
func (videoService *VideoServiceImpl) creatVideo(video *dao.ResVideo, userId string) {
	//防止输出密码
	video.Author.Password = ""
	//建立协程组，当这一组的携程全部完成后，才会结束本方法
	var wg sync.WaitGroup
	wg.Add(5)
	var err error

	u := GetVideoService()
	//插入点赞数量
	go func() {
		video.FavoriteCount, err = u.FavoriteCount(strconv.FormatInt(video.Id, 10), true)
		if err != nil {
			log.Printf("get favorite count failed：%v", err)
		}
		wg.Done()
	}()

	//插入收藏数量
	go func() {
		video.CollectionCount, err = u.CollectionCount(strconv.FormatInt(video.Id, 10), true)
		if err != nil {
			log.Printf("get collection count failed：%v", err)
		}
		wg.Done()
	}()

	//获取该视屏的评论数字
	go func() {
		video.CommentCount, err = u.CountFromVideoId(video.Id)
		if err != nil {
			log.Printf("get comment count failed：%v", err)
		}
		wg.Done()
	}()

	//获取当前用户是否点赞了该视频
	go func() {
		video.IsFavorite, err = u.IsFavorite(strconv.FormatInt(video.Id, 10), userId)
		if err != nil {
			log.Printf("get IsFavorite failed：%v", err)
		}
		wg.Done()
	}()

	//获取当前用户是否点赞了该视频
	go func() {
		video.IsCollection, err = u.IsCollection(strconv.FormatInt(video.Id, 10), userId)
		if err != nil {
			log.Printf("get IsCollection failed：%v", err)
		}
		wg.Done()
	}()

	wg.Wait()
}

// GetVideoService 拼装videoService
func GetVideoService() VideoServiceImpl {
	var userService UserServiceImpl
	var videoService VideoServiceImpl
	var favoriteService FavoriteServiceImpl
	var collectionService CollectionServiceImpl
	var commentService CommentServiceImpl
	userService.FavoriteService = &favoriteService
	userService.CollectionService = &collectionService
	favoriteService.VideoService = &videoService
	collectionService.VideoService = &videoService
	commentService.UserService = &userService
	videoService.CommentService = &commentService
	videoService.FavoriteService = &favoriteService
	videoService.CollectionService = &collectionService
	videoService.UserService = &userService
	return videoService
}

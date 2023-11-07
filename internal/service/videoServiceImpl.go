package service

import (
	"context"
	"fmt"
	"github.com/IRONICBo/QiYin_BE/internal/middleware/rabbitmq"
	requestparams "github.com/IRONICBo/QiYin_BE/internal/params/request"
	"log"
	"strconv"
	"strings"
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

func (videoService *VideoServiceImpl) addHisVideoList(videoId int64, curId string, videoList *[]dao.ResVideo, wg *sync.WaitGroup) (*[]dao.ResVideo, error) {
	defer wg.Done()
	//调用videoService接口，GetVideo：根据videoId，当前用户id:curId，返回Video类型对象
	video, err := videoService.GetVideo(videoId, curId)
	if err != nil {
		log.Println("this history video is miss")
		return videoList, err
	}
	//将Video类型对象添加到集合中去
	*videoList = append(*videoList, video)
	return videoList, nil
}

func (videoService *VideoServiceImpl) GetHisVideos(userId string) ([]dao.ResVideo, error) {
	ctx := context.Background()
	//将int64 userId转换为 string key
	key := fmt.Sprintf("%s:%s", utils.VideoHis, userId)
	//step1:查询Redis CollectionUserId,如果key：strUserId存在,则获取集合中全部videoId
	if n, err := db.GetRedis().Exists(ctx, key).Result(); n > 0 {
		if err != nil {
			log.Printf("history query failed：%v", err)
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
		videoList := new([]dao.ResVideo)
		//采用协程并发将Video类型对象添加到集合中去
		i := len(videoIdList) - 1 //去掉DefaultRedisValue
		if i == 0 {
			return *videoList, nil
		}
		var wg sync.WaitGroup
		wg.Add(i)
		for j := 0; j <= i; j++ {
			//将string videoId转换为 int64 VideoId
			videoId, err := strconv.ParseInt(videoIdList[j], 10, 64)
			if videoId == utils.DefaultRedisValue || err != nil {
				continue
			}
			go videoService.addHisVideoList(videoId, userId, videoList, &wg)
		}
		wg.Wait()
		fmt.Println("后来的", videoList)
		return *videoList, nil
	} else {
		ids, _, err := videoService.videoToRedis(ctx, key, userId)
		if err != nil {
			return nil, err
		}
		//提前开辟点赞列表空间
		videoList := new([]dao.ResVideo)
		i := len(ids)
		if i == 0 {
			return *videoList, nil
		}
		var wg sync.WaitGroup
		wg.Add(i)
		for j := 0; j < i; j++ {
			videoId, _ := strconv.ParseInt(ids[j], 10, 64)
			go videoService.addHisVideoList(videoId, userId, videoList, &wg)
		}
		wg.Wait()
		return *videoList, nil
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

func (videoService *VideoServiceImpl) SaveVideoHis(userId string, param *requestparams.VideoHisParams) error {
	sb := strings.Builder{}
	sb.WriteString(userId)
	sb.WriteString(" ")
	sb.WriteString(strconv.FormatInt(param.VideoId, 10))
	sb.WriteString(" ")
	sb.WriteString(strconv.FormatFloat(param.WatchRatio, 'E', -1, 64))

	ctx := context.Background()
	//将int64 userId转换为 string key
	key := fmt.Sprintf("%s:%s", utils.VideoHis, userId)
	//查询Redis CollectionUserId(key:key)是否已经加载过此信息
	if n, err := db.GetRedis().Exists(ctx, key).Result(); n > 0 {
		//如果有问题，说明查询redis失败,返回错误信息
		if err != nil {
			log.Printf("save history failed：%v", err)
			return err
		} //如果加载过此信息key:key，则加入value:videoId
		if _, err1 := db.GetRedis().SAdd(ctx, key, strconv.FormatInt(param.VideoId, 10)).Result(); err1 != nil {
			log.Print(err1)
			return err1
		} else {
			rabbitmq.RmqVideoAdd.Publish(sb.String())
		}
	} else {
		_, _, err := videoService.videoToRedis(ctx, key, userId)
		if err != nil {
			return err
		}
		if _, err2 := db.GetRedis().SAdd(ctx, key, strconv.FormatInt(param.VideoId, 10)).Result(); err2 != nil {
			return err2
		} else {
			rabbitmq.RmqVideoAdd.Publish(sb.String())
		}
	}

	return nil
}

// 同步点赞状态mysql到redis
func (videoService *VideoServiceImpl) videoToRedis(ctx context.Context, key string, userId string) ([]string, bool, error) {
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
	IdList, err = dao.GeVideoHisList(userId)
	if err != nil {
		return []string{}, false, err
	}

	//维护Redis CollectionUserId(key:key)，遍历videoIdList加入
	for _, id := range IdList {
		if _, err1 := db.GetRedis().SAdd(ctx, key, id).Result(); err1 != nil {
			db.GetRedis().Del(ctx, key)
			log.Println("async failed")
			return []string{}, false, err1
		}
	}
	return IdList, true, nil
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

	var resVideos []dao.ResVideo
	// 得到点赞数和收藏数
	for _, video := range videoList {
		res, err := videoService.creatVideo(&video, userId)
		if err != nil {
			resVideos = append(resVideos, video)
		}
		resVideos = append(resVideos, *res)
	}
	return resVideos, nil
}

// Search 搜索.
func (videoService *VideoServiceImpl) SearchTag(category string, userId string) ([]dao.ResVideo, error) {
	// 查询到相关的videolist + 相关的用户信息
	videoList, err := dao.GetVideoByCate(category)
	// 查询失败直接返回
	if err != nil {
		log.Printf("Search query failed：%v", err)
		return videoList, err
	}

	var resVideos []dao.ResVideo
	// 得到点赞数和收藏数
	for _, video := range videoList {
		res, err := videoService.creatVideo(&video, userId)
		if err != nil {
			resVideos = append(resVideos, video)
		}
		resVideos = append(resVideos, *res)
	}
	return resVideos, nil
}

func (videoService *VideoServiceImpl) GetVideoByUserId(userId string, curUsrId string) ([]dao.ResVideo, error) {
	// 查询到相关的videolist + 相关的用户信息
	videoList, err := dao.GetVideoBuUserId(userId)
	// 查询失败直接返回
	if err != nil {
		log.Printf("query failed：%v", err)
		return videoList, err
	}

	var resVideos []dao.ResVideo
	// 得到点赞数和收藏数
	for _, video := range videoList {
		res, err := videoService.creatVideo(&video, curUsrId)
		if err != nil {
			resVideos = append(resVideos, video)
		}
		resVideos = append(resVideos, *res)
	}
	return resVideos, nil
}

func (videoService *VideoServiceImpl) GetVideoById(videoId string, curUsrId string) (dao.ResVideo, error) {
	// 查询到相关的videolist + 相关的用户信息
	video, err := dao.GetVideo(videoId)
	// 查询失败直接返回
	if err != nil {
		log.Printf("query failed：%v", err)
		return video, err
	}

	// 得到点赞数和收藏数
	res, err := videoService.creatVideo(&video, curUsrId)

	return *res, nil
}

func (videoService *VideoServiceImpl) GetVideos(curUsrId string) ([]dao.ResVideo, error) {
	// 查询到相关的videolist + 相关的用户信息
	videoList, err := dao.GetVideos()
	// 查询失败直接返回
	if err != nil {
		log.Printf("query failed：%v", err)
		return videoList, err
	}

	//var resVideos []dao.ResVideo
	//// 得到点赞数和收藏数
	//for _, video := range videoList {
	//	res, err := videoService.creatVideo(&video, curUsrId)
	//	if err != nil {
	//		resVideos = append(resVideos, video)
	//	}
	//	resVideos = append(resVideos, *res)
	//}
	return videoList, nil
}

func (videoService *VideoServiceImpl) UploadVideo(userId string, param *requestparams.VideoUpdateParams) error {

	newVideo := dao.Video{
		UserId:      userId,
		PlayUrl:     param.PlayUrl,
		CoverUrl:    param.CoverUrl,
		PublishTime: time.Now(),
		Title:       param.Title,
		Desc:        param.Desc,
		Category:    param.Category,
		Tags:        param.Tags,
	}
	return dao.InsertVideo(&newVideo)
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
func (videoService *VideoServiceImpl) creatVideo(video *dao.ResVideo, userId string) (*dao.ResVideo, error) {
	//防止输出密码
	video.Author.Password = ""
	//建立协程组，当这一组的携程全部完成后，才会结束本方法
	var wg sync.WaitGroup

	numOfWait := 3
	if userId != "" {
		numOfWait = 5
	}
	wg.Add(numOfWait)

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

	if userId != "" {
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
	}
	wg.Wait()
	return video, err
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

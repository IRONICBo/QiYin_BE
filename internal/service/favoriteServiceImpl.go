package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/IRONICBo/QiYin_BE/internal/conn/db"
	"github.com/IRONICBo/QiYin_BE/internal/dal/dao"
	"github.com/IRONICBo/QiYin_BE/internal/middleware/rabbitmq"
	"github.com/IRONICBo/QiYin_BE/internal/utils"
	"github.com/gin-gonic/gin"
)

type FavoriteServiceImpl struct {
	VideoService
	UserService
}

// NewFavoriteService return new service with gin context.
func NewFavoriteService(c *gin.Context) *FavoriteServiceImpl {
	return &FavoriteServiceImpl{
		UserService: &UserServiceImpl{
			Service: Service{
				ctx: c,
			},
		},
		VideoService: &VideoServiceImpl{},
	}
}

// IsFavourite 根据favorite:userId,videoId查询点赞状态;.
func (favorite *FavoriteServiceImpl) IsFavourite(videoId string, userId string) (bool, error) {
	key := fmt.Sprintf("%s:%s", utils.Favorite, userId)
	ctx := context.Background()
	if n, err := db.GetRedis().Exists(ctx, key).Result(); n > 1 {
		if err != nil {
			log.Printf("this user has no favorite：%v", err)
			return false, err
		}
		exist, err1 := db.GetRedis().SIsMember(ctx, key, videoId).Result()
		if err1 != nil {
			log.Printf("This user didn't favorite this video：%v", err1)
			return false, err1
		}
		log.Printf("IsFavourite query success")
		return exist, nil
	} else {
		// 反过来查一下
		reverseKey := fmt.Sprintf("%s:%s", utils.Favorite, videoId)
		if n, err := db.GetRedis().Exists(ctx, reverseKey).Result(); n > 1 {
			if err != nil {
				log.Printf("this video has no favorite：%v", err)
				return false, err
			}
			exist, err1 := db.GetRedis().SIsMember(ctx, reverseKey, userId).Result()
			if err1 != nil {
				log.Printf("This video didn't favorite this user：%v", err1)
				return false, err1
			}
			log.Printf("IsFavourite query success")
			return exist, nil
		} else {
			// 两种方式都没有查到  所以都需要进行同步
			ok, err := favorite.favoriteToRedis(ctx, key, userId, false)
			if err != nil || !ok {
				return false, err
			}
			ok, err = favorite.favoriteToRedis(ctx, fmt.Sprintf("%s:%s", utils.Favorite, videoId), videoId, true)
			if err != nil || !ok {
				return false, err
			}
			// 查询Redis FavoriteUserId,key：strUserId中是否存在value:videoId,存在返回true,不存在返回false
			exist, err2 := db.GetRedis().SIsMember(ctx, key, videoId).Result()
			// 如果有问题，说明操作redis失败,返回默认false,返回错误信息
			if err2 != nil {
				return false, err2
			}
			return exist, nil
		}
	}
}

// 同步点赞状态mysql到redis.
func (favorite *FavoriteServiceImpl) favoriteToRedis(ctx context.Context, key string, userId string, isReverse bool) (bool, error) {
	//如果不存在，则维护Redis FavoriteUserId 新建key:key,设置过期时间，加入DefaultRedisValue，
	//key:key，加入value:DefaultRedisValue,过期才会删，防止删最后一个数据的时候数据库还没更新完出现脏读，或者数据库操作失败造成的脏读
	//通过userId查询favorites表,返回所有点赞videoId，加入key:strUserId集合中,
	//再加入当前videoId,再更新favorites表此条数据
	if _, err := db.GetRedis().SAdd(ctx, key, string(rune(utils.DefaultRedisValue))).Result(); err != nil {
		log.Print(err)
		db.GetRedis().Del(ctx, key)
		return false, err
	}
	_, err := db.GetRedis().Expire(ctx, key, time.Duration(utils.OneMonth)*time.Second).Result()
	if err != nil {
		log.Printf("set expire failed")
		db.GetRedis().Del(ctx, key)
		return false, err
	}
	var IdList []string
	// video 作为id
	if isReverse {
		videoId, _ := strconv.ParseInt(userId, 10, 64)
		IdList, err = dao.GetFavoriteUserIdList(videoId)
	} else {
		IdList, err = dao.GetFavoriteVideoIdList(userId)
	}

	// videoIdList, err1 := dao.GetFavoriteVideoIdList(userId)
	// 如果有问题，说明查询失败，返回错误信息："get favoriteVideoIdList failed"
	if err != nil {
		return false, err
	}

	// 维护Redis FavoriteUserId(key:key)，遍历videoIdList加入
	for _, favoriteId := range IdList {
		if _, err1 := db.GetRedis().SAdd(ctx, key, favoriteId).Result(); err1 != nil {
			db.GetRedis().Del(ctx, key)
			log.Println("async failed")
			return false, err1
		}
	}
	return true, nil
}

// 点赞.
func (favorite *FavoriteServiceImpl) favouriteDo(key string, videoId string, userId string, sb strings.Builder, isReverse bool) error {
	ctx := context.Background()
	// 查询Redis FavoriteUserId(key:key)是否已经加载过此信息
	if n, err := db.GetRedis().Exists(ctx, key).Result(); n > 0 {
		// 如果有问题，说明查询redis失败,返回错误信息
		if err != nil {
			log.Printf("FavouriteAction query failed：%v", err)
			return err
		} // 如果加载过此信息key:key，则加入value:videoId
		if _, err1 := db.GetRedis().SAdd(ctx, key, videoId).Result(); err1 != nil {
			log.Print(err1)
			return err1
		} else {
			rabbitmq.RmqFavoriteAdd.Publish(sb.String())
		}
	} else {
		_, err := favorite.favoriteToRedis(ctx, key, userId, isReverse)
		if err != nil {
			return err
		}
		if _, err2 := db.GetRedis().SAdd(ctx, key, videoId).Result(); err2 != nil {
			return err2
		} else {
			rabbitmq.RmqFavoriteAdd.Publish(sb.String())
		}
	}

	return nil
}

func (favorite *FavoriteServiceImpl) favouriteCancel(key string, videoId string, userId string, sb strings.Builder, isReverse bool) error {
	ctx := context.Background()
	// 查询Redis FavoriteUserId(key:key)是否已经加载过此信息
	if n, err := db.GetRedis().Exists(ctx, key).Result(); n > 0 {
		// 如果有问题，说明查询redis失败,返回错误信息
		if err != nil {
			return err
		} // 防止出现redis数据不一致情况，当redis删除操作成功，才执行数据库更新操作
		if _, err1 := db.GetRedis().SRem(ctx, key, videoId).Result(); err1 != nil {
			return err1
		} else {
			// 后续数据库的操作，可以在mq里设置若执行数据库更新操作失败，重新消费该信息

			rabbitmq.RmqFavoriteDel.Publish(sb.String())
		}
	} else {
		_, err := favorite.favoriteToRedis(ctx, key, userId, isReverse)
		if err != nil {
			return err
		}

		if _, err2 := db.GetRedis().SRem(ctx, key, videoId).Result(); err2 != nil {
			return err2
		} else {
			rabbitmq.RmqFavoriteDel.Publish(sb.String())
		}
	}
	return nil
}

// FavouriteAction 根据userId，videoId,actionType对视频进行点赞或者取消赞操作;
// step1: 维护Redis FavoriteUserId(key:key),添加或者删除value:videoId,FavoriteVideoId(key:strVideoId),添加或者删除value:userId;
// step2：更新数据库favorites表;.
func (favorite *FavoriteServiceImpl) FavouriteAction(userId string, videoId string, actionType int32) error {
	key := fmt.Sprintf("%s:%s", utils.Favorite, userId)
	// 将要操作数据库favorites表的信息打入消息队列RmqFavoriteAdd或者RmqFavoriteDel
	// 拼接打入信息
	sb := strings.Builder{}
	sb.WriteString(userId)
	sb.WriteString(" ")
	sb.WriteString(videoId)

	// 执行点赞操作维护
	if actionType == utils.IsFavorite {
		err := favorite.favouriteDo(key, videoId, userId, sb, false)
		if err != nil {
			return err
		}
		err = favorite.favouriteDo(fmt.Sprintf("%s:%s", utils.Favorite, videoId), userId, videoId, sb, true)
		if err != nil {
			return err
		}
	} else { // 执行取消赞操作维护
		err := favorite.favouriteCancel(key, videoId, userId, sb, false)
		if err != nil {
			return err
		}
		err = favorite.favouriteCancel(fmt.Sprintf("%s:%s", utils.Favorite, videoId), userId, videoId, sb, true)
		if err != nil {
			return err
		}
	}
	return nil
}

// TotalFavourite 根据userId获取这个用户总共被点赞数量
// 首先需要找到用户的视频  并统计每个视频的被点赞数.
func (favorite *FavoriteServiceImpl) TotalFavourite(userId string) (int64, error) {
	// 根据userId获取这个用户的发布视频列表信息
	videoIdList, err := favorite.GetVideoIdList(userId)
	if err != nil {
		return 0, err
	}
	var sum int64
	// 提前开辟空间,存取每个视频的点赞数
	videoFavoriteCountList := new([]int64)
	// 采用协程并发将对应videoId的点赞数添加到集合中去
	i := len(videoIdList)
	var wg sync.WaitGroup
	wg.Add(i)
	for j := 0; j < i; j++ {
		go favorite.addVideoFavoriteCount(videoIdList[j], videoFavoriteCountList, &wg)
	}
	wg.Wait()
	// 遍历累加，求总被点赞数
	for _, count := range *videoFavoriteCountList {
		sum += count
	}
	return sum, nil
}

// FavouriteCount 根据videoId获取对应点赞数量;
// step1：查询Redis FavoriteVideoId(key:strVideoId)是否已经加载过此信息，通过set集合中userId个数，获取点赞数量;
// step2：FavoriteVideoId中都没有对应key，维护FavoriteVideoId对应key，再通过set集合中userId个数，获取点赞数量;.
func (favorite *FavoriteServiceImpl) FavouriteCount(videoId string, isReverse bool) (int64, error) {
	// 将int64 videoId转换为 string strVideoId
	key := fmt.Sprintf("%s:%s", utils.Favorite, videoId)
	ctx := context.Background()

	// step1 如果key:strVideoId存在 则计算集合中userId个数   set中会有默认值  所以必须要大于1
	if n, err := db.GetRedis().Exists(ctx, key).Result(); n > 1 {
		// 如果有问题，说明查询redis失败,返回默认false,返回错误信息
		if err != nil {
			log.Printf("FavouriteCount query failed：%v", err)
			return 0, err
		}
		// 获取集合中userId个数
		count, err1 := db.GetRedis().SCard(ctx, videoId).Result()
		// 如果有问题，说明操作redis失败,返回默认0,返回错误信息
		if err1 != nil {
			return 0, err1
		}
		return count - 1, nil // 去掉默认的
	} else {
		_, err1 := favorite.favoriteToRedis(ctx, key, videoId, isReverse)
		if err1 != nil {
			return 0, err1
		}
		// 再通过set集合中userId个数,获取点赞数量
		count, err2 := db.GetRedis().SCard(ctx, key).Result()
		// fmt.Println(key)
		// fmt.Println(db.GetRedis().SMembers(ctx, key))
		// fmt.Println(count)
		// 如果有问题，说明操作redis失败,返回默认0,返回错误信息
		if err2 != nil {
			log.Printf("FavouriteCount query count failed%v", err2)
			return 0, err2
		}
		return count - 1, nil
	}
}

// FavouriteVideoCount 根据userId获取这个用户点赞视频数量.
func (favorite *FavoriteServiceImpl) FavouriteVideoCount(userId string) (int64, error) {
	count, err := favorite.FavouriteCount(userId, false)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// addVideoFavoriteCount 根据videoId，将该视频点赞数加入对应提前开辟好的空间内.
func (favorite *FavoriteServiceImpl) addVideoFavoriteCount(videoId int64, videoFavoriteCountList *[]int64, wg *sync.WaitGroup) {
	defer wg.Done()
	// 调用FavouriteCount：根据videoId,获取点赞数
	count, err := favorite.FavouriteCount(strconv.FormatInt(videoId, 10), true)
	if err != nil {
		// 如果有错误，输出错误信息，并不加入该视频点赞数
		log.Printf(err.Error())
		return
	}
	*videoFavoriteCountList = append(*videoFavoriteCountList, count)
}

//
//// GetFavouriteList 根据userId，curId(当前用户Id),返回userId的点赞列表;
//// step1：查询Redis FavoriteUserId(key:key)是否已经加载过此信息，获取集合中全部videoId，并添加到点赞列表集合中;
//// step2：FavoriteUserId中都没有对应key，维护FavoriteUserId对应key，同时添加到点赞列表集合中;
//func (favorite *FavoriteServiceImpl) GetFavouriteList(userId string, curId string) ([]dao.Video, error) {
//	ctx := context.Background()
//	//将int64 userId转换为 string key
//	key := strconv.FormatInt(userId, 10)
//	//step1:查询Redis FavoriteUserId,如果key：strUserId存在,则获取集合中全部videoId
//	if n, err := db.GetRedis().Exists(ctx, key).Result(); n > 0 {
//		//如果有问题，说明查询redis失败,返回默认nil,返回错误信息
//		if err != nil {
//			log.Printf("方法:GetFavouriteList RedisFavoriteVideoId query key失败：%v", err)
//			return nil, err
//		}
//		//获取集合中全部videoId
//		videoIdList, err1 := db.GetRedis().SMembers(ctx, key).Result()
//		//如果有问题，说明查询redis失败,返回默认nil,返回错误信息
//		if err1 != nil {
//			log.Printf("方法:GetFavouriteList RedisFavoriteVideoId get values失败：%v", err1)
//			return nil, err1
//		}
//		//提前开辟点赞列表空间
//		favoriteVideoList := new([]Video)
//		//采用协程并发将Video类型对象添加到集合中去
//		i := len(videoIdList) - 1 //去掉DefaultRedisValue
//		if i == 0 {
//			return *favoriteVideoList, nil
//		}
//		var wg sync.WaitGroup
//		wg.Add(i)
//		for j := 0; j <= i; j++ {
//			//将string videoId转换为 int64 VideoId
//			videoId, _ := strconv.ParseInt(videoIdList[j], 10, 64)
//			if videoId == config.DefaultRedisValue {
//				continue
//			}
//			go favorite.addFavouriteVideoList(videoId, curId, favoriteVideoList, &wg)
//		}
//		wg.Wait()
//		return *favoriteVideoList, nil
//	} else { //如果Redis FavoriteUserId不存在此key,通过userId查询favorites表,返回所有点赞videoId，并维护到Redis FavoriteUserId(key:key)
//		//key:key，加入value:DefaultRedisValue,过期才会删，防止删最后一个数据的时候数据库还没更新完出现脏读，或者数据库操作失败造成的脏读
//		if _, err := db.GetRedis().SAdd(ctx, key, config.DefaultRedisValue).Result(); err != nil {
//			log.Printf("方法:GetFavouriteList RedisFavoriteUserId add value失败")
//			db.GetRedis().Del(ctx, key)
//			return nil, err
//		}
//		//给键值设置有效期，类似于gc机制
//		_, err := db.GetRedis().Expire(ctx, key,
//			time.Duration(config.OneMonth)*time.Second).Result()
//		if err != nil {
//			log.Printf("方法:GetFavouriteList RedisFavoriteUserId 设置有效期失败")
//			db.GetRedis().Del(ctx, key)
//			return nil, err
//		}
//		videoIdList, err1 := dao.GetFavoriteVideoIdList(userId)
//		//如果有问题，说明查询数据库失败，返回nil和错误信息:"get favoriteVideoIdList failed"
//		if err1 != nil {
//			log.Println(err1.Error())
//			db.GetRedis().Del(ctx, key)
//			return nil, err1
//		}
//		//遍历videoIdList,添加进key的集合中，若失败，删除key，并返回错误信息，这么做的原因是防止脏读，
//		//保证redis与mysql数据一致性
//		for _, favoriteVideoId := range videoIdList {
//			if _, err2 := db.GetRedis().SAdd(ctx, key, favoriteVideoId).Result(); err2 != nil {
//				log.Printf("方法:GetFavouriteList RedisFavoriteUserId add value失败")
//				db.GetRedis().Del(ctx, key)
//				return nil, err2
//			}
//		}
//		//提前开辟点赞列表空间
//		favoriteVideoList := new([]Video)
//		//采用协程并发将Video类型对象添加到集合中去
//		i := len(videoIdList) - 1 //去掉DefaultRedisValue
//		if i == 0 {
//			return *favoriteVideoList, nil
//		}
//		var wg sync.WaitGroup
//		wg.Add(i)
//		for j := 0; j <= i; j++ {
//			if videoIdList[j] == config.DefaultRedisValue {
//				continue
//			}
//			go favorite.addFavouriteVideoList(videoIdList[j], curId, favoriteVideoList, &wg)
//		}
//		wg.Wait()
//		return *favoriteVideoList, nil
//	}
//}

//// addFavouriteVideoList 根据videoId,登录用户curId，添加视频对象到点赞列表空间
//func (favorite *FavoriteServiceImpl) addFavouriteVideoList(videoId int64, curId int64, favoriteVideoList *[]Video, wg *sync.WaitGroup) {
//	defer wg.Done()
//	//调用videoService接口，GetVideo：根据videoId，当前用户id:curId，返回Video类型对象
//	video, err := favorite.GetVideo(videoId, curId)
//	//如果没有获取这个video_id的视频，视频可能被删除了,打印异常,并且不加入此视频
//	if err != nil {
//		log.Println(errors.New("this favourite video is miss"))
//		return
//	}
//	//将Video类型对象添加到集合中去
//	*favoriteVideoList = append(*favoriteVideoList, video)
//}

// GetLikeService 解决likeService调videoService,videoService调userService,useService调likeService循环依赖的问题.
func GetLikeService() FavoriteServiceImpl {
	var userService UserServiceImpl
	var videoService VideoServiceImpl
	var likeService FavoriteServiceImpl
	userService.FavoriteService = &likeService
	likeService.VideoService = &videoService
	videoService.UserService = &userService
	return likeService
}

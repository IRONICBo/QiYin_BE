package dao

import (
	"errors"
	"log"

	"github.com/IRONICBo/QiYin_BE/internal/conn/db"
	"github.com/IRONICBo/QiYin_BE/internal/utils"
)

// Favorite 表的结构。
type Favorite struct {
	Id      int64  // 自增主键
	UserId  string // 点赞用户id
	VideoId int64  // 视频id
	Cancel  int8   // 是否点赞，0为点赞，1为取消赞
}

// TableName 修改表名映射.
func (Favorite) TableName() string {
	return "favorites"
}

// GetFavoriteUserIdList 根据videoId获取点赞userId.
func GetFavoriteUserIdList(videoId int64) ([]string, error) {
	var favoriteUserIdList []string // 存所有该视频点赞用户id；
	// 查询favorites表对应视频id点赞用户，返回查询结果
	err := db.GetMysqlDB().Model(Favorite{}).Where(map[string]interface{}{"video_id": videoId, "cancel": utils.IsFavorite}).
		Pluck("user_id", &favoriteUserIdList).Error
	// 查询过程出现错误，返回默认值0，并输出错误信息
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("get favoriteUserIdList failed")
	} else {
		// 没查询到或者查询到结果，返回数量以及无报错
		return favoriteUserIdList, nil
	}
}

// UpdateFavorite 根据userId，videoId,actionType点赞或者取消赞.
func UpdateFavorite(userId string, videoId int64, actionType int32) error {
	// 更新当前用户观看视频的点赞状态“cancel”，返回错误结果
	err := db.GetMysqlDB().Model(Favorite{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).
		Update("cancel", actionType).Error
	// 如果出现错误，返回更新数据库失败
	if err != nil {
		log.Println(err.Error())
		return errors.New("update data fail")
	}
	// 更新操作成功
	return nil
}

// InsertFavorite 插入点赞数据.
func InsertFavorite(favoriteData Favorite) error {
	// 创建点赞数据，默认为点赞，cancel为0，返回错误结果
	err := db.GetMysqlDB().Model(Favorite{}).Create(&favoriteData).Error
	// 如果有错误结果，返回插入失败
	if err != nil {
		log.Println(err.Error())
		return errors.New("insert data fail")
	}
	return nil
}

// GetFavoriteInfo 根据userId,videoId查询点赞信息.
func GetFavoriteInfo(userId string, videoId int64) (Favorite, error) {
	// 创建一条空favorite结构体，用来存储查询到的信息
	var favoriteInfo Favorite
	// 根据userid,videoId查询是否有该条信息，如果有，存储在favoriteInfo,返回查询结果
	err := db.GetMysqlDB().Model(Favorite{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).
		First(&favoriteInfo).Error
	if err != nil {
		// 查询数据为0，打印"can't find data"，返回空结构体，这时候就应该要考虑是否插入这条数据了
		if "record not found" == err.Error() {
			log.Println("can't find data")
			return Favorite{}, nil
		} else {
			// 如果查询数据库失败，返回获取favoriteInfo信息失败
			log.Println(err.Error())
			return favoriteInfo, errors.New("get favoriteInfo failed")
		}
	}
	return favoriteInfo, nil
}

// GetFavoriteVideoIdList 根据userId查询所属点赞全部videoId.
func GetFavoriteVideoIdList(userId string) ([]string, error) {
	var favoriteVideoIdList []string
	err := db.GetMysqlDB().Model(Favorite{}).Where(map[string]interface{}{"user_id": userId, "cancel": utils.IsFavorite}).
		Pluck("video_id", &favoriteVideoIdList).Error
	if err != nil {
		// 查询数据为0，返回空favoriteVideoIdList切片，以及返回无错误
		if "record not found" == err.Error() {
			log.Println("there are no favoriteVideoId")
			return favoriteVideoIdList, nil
		} else {
			// 如果查询数据库失败，返回获取favoriteVideoIdList失败
			log.Println(err.Error())
			return favoriteVideoIdList, errors.New("get favoriteVideoIdList failed")
		}
	}
	return favoriteVideoIdList, nil
}

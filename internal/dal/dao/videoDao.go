package dao

import (
	"errors"
	"log"
	"time"

	"github.com/IRONICBo/QiYin_BE/internal/conn/db"
)

type Video struct {
	Id          int64     `json:"id"`
	UserId      string    `json:"user_id"`
	PlayUrl     string    `json:"play_url"`
	CoverUrl    string    `json:"cover_url"`
	PublishTime time.Time `json:"publish_time"`
	Title       string    `json:"title"` // 视频名
	Desc        string    `json:"desc"`
	Category    int64     `json:"category"`
	Tags        string    `json:"tags"`
}

type ResVideo struct {
	Video
	Author          User  `json:"author" gorm:"foreignKey:UserId"`
	FavoriteCount   int64 `json:"favorite_count"`
	CommentCount    int64 `json:"comment_count"`
	IsFavorite      bool  `json:"is_favorite"`
	CollectionCount int64 `json:"collection_count"`
	IsCollection    bool  `json:"is_collection"`
}

// 浏览记录
type UserWatchAction struct {
	Id         int64     `json:"id"`
	UserId     string    `json:"user_id"`
	VideoId    int64     `json:"video_id"`    // 视频ID
	WatchRatio float64   `json:"watch_ratio"` // 观看完成比例
	FinishTime time.Time `json:"finish_time"` // 用户Action时间戳
}

// GetVideoIdsByAuthorId
// 通过作者id来查询发布的视频id切片集合.
func GetVideoIdsByAuthorId(userId string) ([]int64, error) {
	var id []int64
	// 通过pluck来获得单独的切片
	result := db.GetMysqlDB().Model(&Video{}).Where("user_id", userId).Pluck("id", &id)
	// 如果出现问题，返回对应到空，并且返回error
	if result.Error != nil {
		return nil, result.Error
	}
	return id, nil
}

// GetVideoByTitle
// 通过关键字搜索视频  title.
func GetVideoByTitle(value string) ([]ResVideo, error) {
	var videoList []ResVideo
	result := db.GetMysqlDB().Table("videos").Where("title LIKE ?", "%"+value+"%").Preload("Author").Order("publish_time desc").Find(&videoList)
	// 如果出现问题，返回对应到空，并且返回error
	if result.Error != nil {
		return []ResVideo{}, result.Error
	}
	return videoList, nil
}

func GetVideoByCate(value string) ([]ResVideo, error) {
	var videoList []ResVideo
	result := db.GetMysqlDB().Table("videos").Where("category = ?", value).Preload("Author").Order("publish_time desc").Find(&videoList)
	// 如果出现问题，返回对应到空，并且返回error
	if result.Error != nil {
		return []ResVideo{}, result.Error
	}
	return videoList, nil
}

// GetVideoBuUserId
func GetVideoBuUserId(value string) ([]ResVideo, error) {
	var video []ResVideo
	result := db.GetMysqlDB().Table("videos").Where("user_id = ?", value).Preload("Author").Order("publish_time desc").Find(&video)
	// 如果出现问题，返回对应到空，并且返回error
	if result.Error != nil {
		return []ResVideo{}, result.Error
	}
	return video, nil
}

// GetVideo
func GetVideo(value string) (ResVideo, error) {
	var video ResVideo
	result := db.GetMysqlDB().Table("videos").Where("id = ?", value).Preload("Author").Order("publish_time desc").First(&video)
	// 如果出现问题，返回对应到空，并且返回error
	if result.Error != nil {
		return ResVideo{}, result.Error
	}
	return video, nil
}

// GetVideos
func GetVideos() ([]ResVideo, error) {
	var videoList []ResVideo
	result := db.GetMysqlDB().Table("videos").Preload("Author").Order("publish_time desc").Find(&videoList)
	// 如果出现问题，返回对应到空，并且返回error
	if result.Error != nil {
		return []ResVideo{}, result.Error
	}
	return videoList, nil
}

// GetVideoById
// 通过userId 搜索视频.
func GetVideoById(videoId int64) (ResVideo, error) {
	var videoList ResVideo
	err := db.GetMysqlDB().Table("videos").Where("id", videoId).Preload("Author").First(&videoList).Error
	if err != nil {
		// 查询数据为0，返回空collectionVideoIdList切片，以及返回无错误
		if "record not found" == err.Error() {
			return videoList, nil
		} else {
			// 如果查询数据库失败，返回获取collectionVideoIdList失败
			log.Println(err.Error())
		}
	}
	return videoList, nil
}

func InsertVideo(video *Video) error {
	err := db.GetMysqlDB().Create(&video).Error
	return err
}

func GetVideoHis(userId string, videoId int64) (UserWatchAction, error) {
	var videoList UserWatchAction
	res := db.GetMysqlDB().Model(UserWatchAction{}).Where("user_id = ? and video_id = ?", userId, videoId).First(&videoList)
	if res.Error != nil {
		return UserWatchAction{}, res.Error
	}
	return videoList, nil
}

func InsertVideoHis(video *UserWatchAction) error {
	err := db.GetMysqlDB().Model(UserWatchAction{}).Create(&video).Error
	return err
}

func UpdateVideoHis(video *UserWatchAction) error {
	err := db.GetMysqlDB().Model(UserWatchAction{}).Where(map[string]interface{}{"user_id": video.UserId, "video_id": video.VideoId}).
		Updates(video).Error
	// 如果出现错误，返回更新数据库失败
	if err != nil {
		log.Println(err.Error())
		return errors.New("update data fail")
	}
	return nil
}

func GeVideoHisList(userId string) ([]string, error) {
	var IdList []string
	err := db.GetMysqlDB().Model(UserWatchAction{}).Where("user_id = ?", userId).Pluck("video_id", &IdList).Error
	// 如果出现错误，返回更新数据库失败
	if err != nil {
		// 查询数据为0，返回空collectionVideoIdList切片，以及返回无错误
		if "record not found" == err.Error() {
			log.Println("there are no id")
			return IdList, nil
		} else {
			// 如果查询数据库失败，返回获取collectionVideoIdList失败
			log.Println(err.Error())
		}
	}
	return IdList, nil

}

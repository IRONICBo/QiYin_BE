package dao

import (
	"time"

	"github.com/IRONICBo/QiYin_BE/internal/conn/db"
)

type Video struct {
	Id          int64     `json:"id"`
	UserId      string    `json:"user_id"`
	PlayUrl     string    `json:"play_url"`
	CoverUrl    string    `json:"cover_url"`
	PublishTime time.Time `json:"publish_time"`
	Title       string    `json:"title"` //视频名
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
	result := db.GetMysqlDB().Table("videos").Where("title LIKE ?", "%"+value+"%").Preload("Author").Find(&videoList)
	// 如果出现问题，返回对应到空，并且返回error
	if result.Error != nil {
		return []ResVideo{}, result.Error
	}
	return videoList, nil
}

// GetVideoById
// 通过userId 搜索视频
func GetVideoById(videoId int64) (ResVideo, error) {
	var videoList ResVideo
	result := db.GetMysqlDB().Table("videos").Where("id", videoId).Preload("Author").First(&videoList)
	//如果出现问题，返回对应到空，并且返回error
	if result.Error != nil {
		return ResVideo{}, result.Error
	}
	return videoList, nil
}

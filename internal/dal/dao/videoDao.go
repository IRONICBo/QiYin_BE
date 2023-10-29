package dao

import "time"

type Video struct {
	Id          int64 `json:"id"`
	UserId      int64
	PlayUrl     string `json:"play_url"`
	CoverUrl    string `json:"cover_url"`
	PublishTime time.Time
	Title       string `json:"title"` //视频名，5.23添加
}

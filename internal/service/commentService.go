package service

import (
	"github.com/IRONICBo/QiYin_BE/internal/dal/dao"
)

// CommentService 接口定义
// 发表评论-使用的结构体-service层引用dao层↑的Comment。
type CommentService interface {
	CountFromVideoId(videoId int64) (int64, error)
	CommentAdd(comment dao.Comment) (dao.CommentData, error)
	CommentDelete(userId string, commentId int64)
	GetList(videoId string) ([]dao.CommentData, error)
	insertRedisVideoCommentId(videoId string, commentId string)
}

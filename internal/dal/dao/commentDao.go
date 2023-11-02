package dao

import (
	"errors"
	"log"
	"time"

	"github.com/IRONICBo/QiYin_BE/internal/conn/db"
	"github.com/IRONICBo/QiYin_BE/internal/utils"
)

// Comment
// 评论信息-数据库中的结构体-dao层使用.
type Comment struct {
	Id          int64     // 评论id
	UserId      string    // 评论用户id
	VideoId     int64     // 视频id
	CommentText string    // 评论内容
	CreateDate  time.Time // 评论发布的日期mm-dd
	Cancel      int32     // 取消评论为1，发布评论为0
}

type CommentData struct {
	Id         int64   `json:"id,omitempty"`
	UserInfo   ResUser `json:"user_info"`
	Content    string  `json:"content,omitempty"`
	CreateDate string  `json:"create_date,omitempty"`
}

// TableName 修改表名映射.
func (Comment) TableName() string {
	return "comments"
}

// Count
// 使用video id 查询Comment数量.
func Count(videoId int64) (int64, error) {
	log.Println("CommentDao-Count: running") // 函数已运行
	// Init()
	var count int64
	// 数据库中查询评论数量
	err := db.GetMysqlDB().Model(Comment{}).Where(map[string]interface{}{"video_id": videoId, "cancel": utils.ValidComment}).Count(&count).Error
	if err != nil {
		log.Println("CommentDao-Count: return count failed") // 函数返回提示错误信息
		return -1, errors.New("find comments count failed")
	}
	log.Println("CommentDao-Count: return count success") // 函数执行成功，返回正确信息
	return count, nil
}

// CommentIdList 根据视频id获取评论id 列表.
func CommentIdList(videoId int64) ([]string, error) {
	var commentIdList []string
	err := db.GetMysqlDB().Model(Comment{}).Select("id").Where("video_id = ?", videoId).Find(&commentIdList).Error
	if err != nil {
		log.Println("CommentIdList:", err)
		return nil, err
	}
	return commentIdList, nil
}

// InsertComment
// 发表评论.
func InsertComment(comment Comment) (int64, error) {
	// 数据库中插入一条评论信息
	err := db.GetMysqlDB().Model(Comment{}).Create(&comment).Error
	if err != nil {
		log.Println("CommentDao-InsertComment: return create comment failed") // 函数返回提示错误信息
		return 0, errors.New("create comment failed")
	}
	log.Println("CommentDao-InsertComment: return success") // 函数执行成功，返回正确信息
	return comment.Id, nil
}

// DeleteComment.
func DeleteComment(id int64) error {
	// 数据库中删除评论-更新评论状态为-1
	err := db.GetMysqlDB().Model(Comment{}).Where("id = ?", id).Update("cancel", utils.InvalidComment).Error
	if err != nil {
		log.Println("CommentDao-DeleteComment: return del comment failed") // 函数返回提示错误信息
		return errors.New("del comment failed")
	}
	log.Println("CommentDao-DeleteComment: return success") // 函数执行成功，返回正确信息
	return nil
}

// GetCommentList.
func GetCommentList(videoId string) ([]Comment, error) {
	// 数据库中查询评论信息list
	var commentList []Comment
	result := db.GetMysqlDB().Model(Comment{}).Where(map[string]interface{}{"video_id": videoId, "cancel": utils.ValidComment}).
		Order("create_date desc").Find(&commentList)
	// 若此视频没有评论信息，返回空列表，不报错
	if result.RowsAffected == 0 {
		log.Println("CommentDao-GetCommentList: return there are no comments") // 函数返回提示无评论
		return nil, nil
	}
	// 若获取评论列表出错
	if result.Error != nil {
		log.Println(result.Error.Error())
		log.Println("CommentDao-GetCommentList: return get comment list failed") // 函数返回提示获取评论错误
		return commentList, errors.New("get comment list failed")
	}
	log.Println("CommentDao-GetCommentList: return commentList success") // 函数执行成功，返回正确信息
	return commentList, nil
}

func GetUserIdByCommentId(commentId int64, userId string) error {
	var commentInfo Comment
	// 先查询是否有此评论
	result := db.GetMysqlDB().Model(Comment{}).Where(map[string]interface{}{"id": commentId, "cancel": utils.ValidComment}).First(&commentInfo)
	if result.RowsAffected == 0 { // 查询到此评论数量为0则返回无此评论
		return errors.New("del comment is not exist")
	} else if commentInfo.UserId != userId {
		return errors.New("have no Permission")
	}
	return nil
}

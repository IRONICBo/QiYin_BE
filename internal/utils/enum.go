package utils

const (
	// redis中key的状态
	Favorite          = "favorite"
	Colloection       = "collection"
	Comment           = "comment"
	CommentCV         = "commentCV"
	DefaultRedisValue = -1 //防止赃读

	OneMonth = 60 * 60 * 24 * 30
	DateTime = "2006-01-02 15:04:05"

	IsFavorite = 1
	Unlike     = -1

	Attempts = 3 // 最大操作次数

	ValidComment   = 1
	InvalidComment = -1
)

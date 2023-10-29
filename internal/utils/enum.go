package utils

const (
	// redis中key的状态
	Favorite          = "favorite"
	Relation          = "relation"
	DefaultRedisValue = -1 //防止赃读

	OneMonth = 60 * 60 * 24 * 30

	IsFavorite = 1
	Unlike     = -1

	Attempts = 3 // 最大操作次数
)

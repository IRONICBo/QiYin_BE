package service

// CollectionService 定义点赞状态和点赞数量
type CollectionService interface {
	IsCollection(videoId string, userId string) (bool, error)
	CollectionCount(videoId string, isReverse bool) (int64, error)
	TotalCollection(userId string) (int64, error)
	CollectionVideoCount(userId string) (int64, error)

	CollectionAction(userId string, videoId string, actionType int32) error
}

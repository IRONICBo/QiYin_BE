package service

// FavoriteService 定义点赞状态和点赞数量.
type FavoriteService interface {
	/*
	   1.其他模块(video)需要使用的业务方法。
	*/
	//IsFavorite 根据当前视频id判断是否点赞了该视频。
	IsFavorite(videoId string, userId string) (bool, error)
	//FavoriteCount 根据当前视频id获取当前视频点赞数量。
	FavoriteCount(videoId string, isReverse bool) (int64, error)
	//TotalFavorite 根据userId获取这个用户总共被点赞数量
	TotalFavorite(userId string) (int64, error)
	//FavoriteVideoCount 根据userId获取这个用户点赞视频数量
	FavoriteVideoCount(userId string) (int64, error)
	/*
	   2.request需要实现的功能
	*/
	//当前用户对视频的点赞操作 ,并把这个行为更新到favorite表中。
	//当前操作行为，1点赞，2取消点赞。
	FavoriteAction(userId string, videoId string, actionType int32) error
	// GetFavoriteList 获取当前用户的所有点赞视频，调用videoService的方法
	//GetFavoriteList(userId int64, curId int64) ([]Video, error)
}

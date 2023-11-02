package service

import "github.com/IRONICBo/QiYin_BE/internal/dal/dao"

type VideoService interface {
	GetVideoIdList(userId string) ([]int64, error)

	Search(searchValue string, userId string) ([]dao.ResVideo, error)

	GetVideo(videoId int64, curUserId string) (dao.ResVideo, error)
}

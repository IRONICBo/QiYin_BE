package service

import (
	"github.com/IRONICBo/QiYin_BE/internal/dal/dao"
	requestparams "github.com/IRONICBo/QiYin_BE/internal/params/request"
)

type VideoService interface {
	GetVideoIdList(userId string) ([]int64, error)

	Search(searchValue string, userId string) ([]dao.ResVideo, error)

	GetVideo(videoId int64, curUserId string) (dao.ResVideo, error)

	GetHisVideos(userId string) ([]dao.ResVideo, error)
	SaveVideoHis(userId string, param *requestparams.VideoHisParams) error
	UploadVideo(userId string, param *requestparams.VideoUpdateParams) error
	GetVideoByUserId(userId string, curUsrId string) ([]dao.ResVideo, error)
}

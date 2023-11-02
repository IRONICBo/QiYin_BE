package service

import (
	"github.com/IRONICBo/QiYin_BE/internal/config"
	responseparams "github.com/IRONICBo/QiYin_BE/internal/params/response"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

// QiNiuService qiniu service.
type QiNiuService struct {
	Service
}

// NewQiNiuService return new service with gin context.
func NewQiNiuService(c *gin.Context) *QiNiuService {
	return &QiNiuService{
		Service: Service{
			ctx: c,
		},
	}
}

// GetUploadToken get qiniu upload token.
func (svc *QiNiuService) GetUploadToken() responseparams.QiNiuTokenResponse {
	accessKey := config.Config.QiNiu.AccessKey
	secretKey := config.Config.QiNiu.SecretKey
	putPolicy := storage.PutPolicy{
		Scope:      config.Config.QiNiu.Bucket,
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","bucket":"$(bucket)","fsize":$(fsize),"name":"$(x:name)"}`,
	}

	mac := qbox.NewMac(accessKey, secretKey)
	uploadToken := putPolicy.UploadToken(mac)

	token := responseparams.QiNiuTokenResponse{
		UploadToken: uploadToken,
	}

	return token
}

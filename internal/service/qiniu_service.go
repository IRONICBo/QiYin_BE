package service

import (
	"encoding/base64"
	"fmt"

	"github.com/IRONICBo/QiYin_BE/internal/config"
	responseparams "github.com/IRONICBo/QiYin_BE/internal/params/response"
	"github.com/IRONICBo/QiYin_BE/pkg/log"
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
func (svc *QiNiuService) GetUploadToken(ticket string) responseparams.QiNiuTokenResponse {
	accessKey := config.Config.QiNiu.AccessKey
	secretKey := config.Config.QiNiu.SecretKey
	putPolicy := storage.PutPolicy{
		Scope:      config.Config.QiNiu.Bucket,
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","bucket":"$(bucket)","fsize":$(fsize),"name":"$(x:name)"}`,
		// Set persistentOps
		PersistentOps: getFopVSampleCommand(ticket, 3),
	}

	mac := qbox.NewMac(accessKey, secretKey)
	uploadToken := putPolicy.UploadToken(mac)

	token := responseparams.QiNiuTokenResponse{
		UploadToken: uploadToken,
	}

	return token
}

// StartSamplePfpop start pfpop sample pop.
func (svc *QiNiuService) StartSamplePfpop(key string, numFrames int) error {
	accessKey := config.Config.QiNiu.AccessKey
	secretKey := config.Config.QiNiu.SecretKey
	bucket := config.Config.QiNiu.Bucket
	mac := qbox.NewMac(accessKey, secretKey)
	cfg := storage.Config{
		UseHTTPS: false,
	}

	operationManager := storage.NewOperationManager(mac, &cfg)
	fopVsample := fmt.Sprintf("vsample/jpg/frames/%d/pattern/%s",
		numFrames,
		base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s_$(count).jpg", key))),
	)
	force := true

	_, err := operationManager.Pfop(bucket, key, fopVsample, "", "", force)
	if err != nil {
		log.Error("StartSamplePfpop", "Pfop error", err)
		return err
	}

	return nil
}

func getFopVSampleCommand(key string, numFrames int) string {
	fopVsample := fmt.Sprintf("vsample/jpg/frames/%d/pattern/%s",
		numFrames,
		base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s_$(count).jpg", key))),
	)
	return fopVsample
}

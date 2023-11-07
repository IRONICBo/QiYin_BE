package main

import (
	"flag"
	"fmt"

	"github.com/IRONICBo/QiYin_BE/internal/middleware/rabbitmq"

	"github.com/IRONICBo/QiYin_BE/internal/config"
	"github.com/IRONICBo/QiYin_BE/internal/conn/db"
	"github.com/IRONICBo/QiYin_BE/internal/router"
	"github.com/IRONICBo/QiYin_BE/internal/utils"
	"github.com/IRONICBo/QiYin_BE/pkg/log"
	"github.com/IRONICBo/QiYin_BE/pkg/server"
)

func init() {
	// arg
	configPath := flag.String("c", "./config.yaml", "config file path")
	flag.Parse()

	// init
	config.ConfigInit(*configPath)
	utils.Banner()
	log.InitLogger()
	db.InitMysqlDB()
	db.InitRedisDB()
	rabbitmq.InitRabbitMQ()
	rabbitmq.InitFavoriteRabbitMQ()
	rabbitmq.InitCommentRabbitMQ()
	rabbitmq.InitCollectionRabbitMQ()
	rabbitmq.InitVideoRabbitMQ()
}

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

// @title QiYin Backend
// @version v0.0.0
// @description QiYin Backend API Docs.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization.
func main() {
	serverAddress := fmt.Sprintf("%s:%d", config.Config.Server.Ip, config.Config.Server.Port)

	r := router.InitRouter()
	s := server.InitServer(serverAddress, r)
	log.Error("server start error: %v", s.ListenAndServe().Error())
}

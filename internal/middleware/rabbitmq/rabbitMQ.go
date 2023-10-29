package rabbitmq

import (
	"fmt"
	"github.com/IRONICBo/QiYin_BE/internal/config"
	"github.com/IRONICBo/QiYin_BE/pkg/log"
	"github.com/streadway/amqp"
	_ "gorm.io/gorm"
)

type RabbitMQ struct {
	conn  *amqp.Connection
	mqurl string
}

var rmq *RabbitMQ

// InitRabbitMQ init mysql connection.
func InitRabbitMQ() {
	// try to use default database [mysql]
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		config.Config.RabbitMQ.Username,
		config.Config.RabbitMQ.Password,
		config.Config.RabbitMQ.Ip,
		config.Config.RabbitMQ.Port,
	)

	Rmq := &RabbitMQ{
		mqurl: dsn,
	}
	dial, err := amqp.Dial(Rmq.mqurl)
	Rmq.failOnErr(err, "创建连接失败")
	Rmq.conn = dial

	rmq = Rmq
	log.Info("RabbitMQ", "connect ok", dsn)
}

// GetMysqlDB get mysql connection.
func GetRabbitMQ() *RabbitMQ {
	return rmq
}

// 连接出错时，输出错误信息。
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Error("%s:%s\n", err, message)
		panic(fmt.Sprintf("%s:%s\n", err, message))
	}
}

// 关闭mq通道和mq的连接。
func (r *RabbitMQ) destroy() {
	err := r.conn.Close()
	if err != nil {
		return
	}
}

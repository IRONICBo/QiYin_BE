package rabbitmq

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/IRONICBo/QiYin_BE/internal/dal/dao"
	"github.com/IRONICBo/QiYin_BE/internal/utils"
	"github.com/streadway/amqp"
)

type VideoMQ struct {
	RabbitMQ
	channel   *amqp.Channel
	queueName string
	exchange  string
	key       string
}

// NewVideoRabbitMQ 获取videoMQ的对应队列。
func NewVideoRabbitMQ(queueName string) *VideoMQ {
	videoMQ := &VideoMQ{
		RabbitMQ:  *GetRabbitMQ(),
		queueName: queueName,
	}
	cha, err := videoMQ.conn.Channel()
	videoMQ.channel = cha
	GetRabbitMQ().failOnErr(err, "获取通道失败")
	return videoMQ
}

// Publish video操作的发布配置。
func (l *VideoMQ) Publish(message string) {
	_, err := l.channel.QueueDeclare(
		l.queueName,
		// 是否持久化
		false,
		// 是否为自动删除
		false,
		// 是否具有排他性
		false,
		// 是否阻塞
		false,
		// 额外属性
		nil,
	)
	if err != nil {
		panic(err)
	}

	err1 := l.channel.Publish(
		l.exchange,
		l.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err1 != nil {
		panic(err)
	}
}

// Consumer video关系的消费逻辑。
func (l *VideoMQ) Consumer() {
	_, err := l.channel.QueueDeclare(l.queueName, false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	// 2、接收消息
	messages, err1 := l.channel.Consume(
		l.queueName,
		// 用来区分多个消费者
		"",
		// 是否自动应答
		true,
		// 是否具有排他性
		false,
		// 如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		// 消息队列是否阻塞
		false,
		nil,
	)
	if err1 != nil {
		panic(err1)
	}

	forever := make(chan bool)
	go l.consumerVideoAdd(messages)

	log.Printf("[*] Waiting for messagees,To exit press CTRL+C")

	<-forever
}

// consumerVideoAdd 赞关系添加的消费方式。
func (l *VideoMQ) consumerVideoAdd(messages <-chan amqp.Delivery) {
	for d := range messages {
		params := strings.Split(fmt.Sprintf("%s", d.Body), " ")
		userId := params[0]
		videoId, _ := strconv.ParseInt(params[1], 10, 64)
		watchRatio, _ := strconv.ParseFloat(params[2], 64)
		// 最多尝试操作数据库的次数
		for i := 0; i < utils.Attempts; i++ {
			flag := false // 默认无问题
			var videoData dao.UserWatchAction
			videoInfo, _ := dao.GetVideoHis(userId, videoId)
			if videoInfo == (dao.UserWatchAction{}) { // 没查到这条数据，则新建这条数据；
				videoData.UserId = userId         // 插入userId
				videoData.VideoId = videoId       // 插入videoId
				videoData.WatchRatio = watchRatio // 插入观看记录
				videoData.FinishTime = time.Now()
				if err := dao.InsertVideoHis(&videoData); err != nil {
					log.Printf(err.Error())
					flag = true // 出现问题
				}
			} else { // 查到这条数据,更新即可;
				if err := dao.UpdateVideoHis(&videoData); err != nil {
					log.Printf(err.Error())
					flag = true // 出现问题
				}
			}
			if flag == false {
				break
			}
		}
	}
}

var (
	RmqVideoAdd *VideoMQ
)

// InitVideoRabbitMQ 初始化rabbitMQ连接。
func InitVideoRabbitMQ() {
	RmqVideoAdd = NewVideoRabbitMQ("video_add")
	go RmqVideoAdd.Consumer()

}

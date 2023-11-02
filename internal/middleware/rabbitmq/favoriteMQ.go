package rabbitmq

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/IRONICBo/QiYin_BE/internal/dal/dao"
	"github.com/IRONICBo/QiYin_BE/internal/utils"
	"github.com/streadway/amqp"
)

type FavoriteMQ struct {
	RabbitMQ
	channel   *amqp.Channel
	queueName string
	exchange  string
	key       string
}

// NewFavoriteRabbitMQ 获取likeMQ的对应队列。
func NewFavoriteRabbitMQ(queueName string) *FavoriteMQ {
	likeMQ := &FavoriteMQ{
		RabbitMQ:  *GetRabbitMQ(),
		queueName: queueName,
	}
	cha, err := likeMQ.conn.Channel()
	likeMQ.channel = cha
	GetRabbitMQ().failOnErr(err, "获取通道失败")
	return likeMQ
}

// Publish like操作的发布配置。
func (l *FavoriteMQ) Publish(message string) {
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

// Consumer like关系的消费逻辑。
func (l *FavoriteMQ) Consumer() {
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
	switch l.queueName {
	case "favorite_add":
		// 点赞消费队列
		go l.consumerFavoriteAdd(messages)
	case "favorite_del":
		// 取消赞消费队列
		go l.consumerFavoriteDel(messages)
	}

	log.Printf("[*] Waiting for messagees,To exit press CTRL+C")

	<-forever
}

// consumerFavoriteAdd 赞关系添加的消费方式。
func (l *FavoriteMQ) consumerFavoriteAdd(messages <-chan amqp.Delivery) {
	for d := range messages {
		// 参数解析。
		params := strings.Split(fmt.Sprintf("%s", d.Body), " ")
		userId := params[0]
		videoId, _ := strconv.ParseInt(params[1], 10, 64)
		// 最多尝试操作数据库的次数
		for i := 0; i < utils.Attempts; i++ {
			flag := false // 默认无问题
			// 如果查询没有数据，用来生成该条点赞信息，存储在likeData中
			var likeData dao.Favorite
			// 先查询是否有这条数据
			likeInfo, err := dao.GetFavoriteInfo(userId, videoId)
			// 如果有问题，说明查询数据库失败，打印错误信息err:"get likeInfo failed"
			if err != nil {
				log.Printf(err.Error())
				flag = true // 出现问题
			} else {
				if likeInfo == (dao.Favorite{}) { // 没查到这条数据，则新建这条数据；
					likeData.UserId = string(userId)   // 插入userId
					likeData.VideoId = videoId         // 插入videoId
					likeData.Cancel = utils.IsFavorite // 插入点赞cancel=0
					// 如果有问题，说明插入数据库失败，打印错误信息err:"insert data fail"
					if err := dao.InsertFavorite(likeData); err != nil {
						log.Printf(err.Error())
						flag = true // 出现问题
					}
				} else { // 查到这条数据,更新即可;
					// 如果有问题，说明插入数据库失败，打印错误信息err:"update data fail"
					if err := dao.UpdateFavorite(userId, videoId, utils.IsFavorite); err != nil {
						log.Printf(err.Error())
						flag = true // 出现问题
					}
				}
				// 一遍流程下来正常执行了，那就打断结束，不再尝试
				if flag == false {
					break
				}
			}
		}
	}
}

// consumerFavoriteDel 赞关系删除的消费方式。
func (l *FavoriteMQ) consumerFavoriteDel(messages <-chan amqp.Delivery) {
	for d := range messages {
		// 参数解析。
		params := strings.Split(fmt.Sprintf("%s", d.Body), " ")
		userId := params[0]
		videoId, _ := strconv.ParseInt(params[1], 10, 64)
		// 最多尝试操作数据库的次数
		for i := 0; i < utils.Attempts; i++ {
			flag := false // 默认无问题
			// 取消赞行为，只有当前状态是点赞状态才会发起取消赞行为，所以如果查询到，必然是cancel==0(点赞)
			// 先查询是否有这条数据
			likeInfo, err := dao.GetFavoriteInfo(userId, videoId)
			// 如果有问题，说明查询数据库失败，返回错误信息err:"get likeInfo failed"
			if err != nil {
				log.Printf(err.Error())
				flag = true // 出现问题
			} else {
				if likeInfo == (dao.Favorite{}) { // 只有当前是点赞状态才能取消点赞这个行为
					// 所以如果查询不到数据则返回错误信息:"can't find data,this action invalid"
					log.Printf(errors.New("can't find data,this action invalid").Error())
				} else {
					// 如果查询到数据，则更新为取消赞状态
					// 如果有问题，说明插入数据库失败，打印错误信息err:"update data fail"
					if err := dao.UpdateFavorite(userId, videoId, utils.Unlike); err != nil {
						log.Printf(err.Error())
						flag = true
					}
				}
			}
			// 一遍流程下来正常执行了，那就打断结束，不再尝试
			if flag == false {
				break
			}
		}
	}
}

var (
	RmqFavoriteAdd *FavoriteMQ
	RmqFavoriteDel *FavoriteMQ
)

// InitFavoriteRabbitMQ 初始化rabbitMQ连接。
func InitFavoriteRabbitMQ() {
	RmqFavoriteAdd = NewFavoriteRabbitMQ("favorite_add")
	go RmqFavoriteAdd.Consumer()

	RmqFavoriteDel = NewFavoriteRabbitMQ("favorite_del")
	go RmqFavoriteDel.Consumer()
}

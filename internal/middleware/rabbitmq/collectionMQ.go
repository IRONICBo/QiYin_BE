package rabbitmq

import (
	"errors"
	"fmt"
	"github.com/IRONICBo/QiYin_BE/internal/dal/dao"
	"github.com/IRONICBo/QiYin_BE/internal/utils"
	"github.com/streadway/amqp"
	"log"
	"strconv"
	"strings"
)

type CollectionMQ struct {
	RabbitMQ
	channel   *amqp.Channel
	queueName string
	exchange  string
	key       string
}

// NewCollectionRabbitMQ 获取collectionMQ的对应队列。
func NewCollectionRabbitMQ(queueName string) *CollectionMQ {
	collectionMQ := &CollectionMQ{
		RabbitMQ:  *GetRabbitMQ(),
		queueName: queueName,
	}
	cha, err := collectionMQ.conn.Channel()
	collectionMQ.channel = cha
	GetRabbitMQ().failOnErr(err, "获取通道失败")
	return collectionMQ
}

// Publish collection操作的发布配置。
func (l *CollectionMQ) Publish(message string) {

	_, err := l.channel.QueueDeclare(
		l.queueName,
		//是否持久化
		false,
		//是否为自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		//额外属性
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

// Consumer collection关系的消费逻辑。
func (l *CollectionMQ) Consumer() {

	_, err := l.channel.QueueDeclare(l.queueName, false, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	//2、接收消息
	messages, err1 := l.channel.Consume(
		l.queueName,
		//用来区分多个消费者
		"",
		//是否自动应答
		true,
		//是否具有排他性
		false,
		//如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		//消息队列是否阻塞
		false,
		nil,
	)
	if err1 != nil {
		panic(err1)
	}

	forever := make(chan bool)
	switch l.queueName {
	case "collection_add":
		//点赞消费队列
		go l.consumerCollectionAdd(messages)
	case "collection_del":
		//取消赞消费队列
		go l.consumerCollectionDel(messages)

	}

	log.Printf("[*] Waiting for messagees,To exit press CTRL+C")

	<-forever

}

// consumerCollectionAdd 赞关系添加的消费方式。
func (l *CollectionMQ) consumerCollectionAdd(messages <-chan amqp.Delivery) {
	for d := range messages {
		// 参数解析。
		params := strings.Split(fmt.Sprintf("%s", d.Body), " ")
		userId := params[0]
		videoId, _ := strconv.ParseInt(params[1], 10, 64)
		//最多尝试操作数据库的次数
		for i := 0; i < utils.Attempts; i++ {
			flag := false //默认无问题
			//如果查询没有数据，用来生成该条点赞信息，存储在collectionData中
			var collectionData dao.Collection
			//先查询是否有这条数据
			collectionInfo, err := dao.GetCollectionInfo(userId, videoId)
			//如果有问题，说明查询数据库失败，打印错误信息err:"get collectionInfo failed"
			if err != nil {
				log.Printf(err.Error())
				flag = true //出现问题
			} else {
				if collectionInfo == (dao.Collection{}) { //没查到这条数据，则新建这条数据；
					collectionData.UserId = string(userId)     //插入userId
					collectionData.VideoId = videoId           //插入videoId
					collectionData.Cancel = utils.IsCollection //插入点赞cancel=0
					//如果有问题，说明插入数据库失败，打印错误信息err:"insert data fail"
					if err := dao.InsertCollection(collectionData); err != nil {
						log.Printf(err.Error())
						flag = true //出现问题
					}
				} else { //查到这条数据,更新即可;
					//如果有问题，说明插入数据库失败，打印错误信息err:"update data fail"
					if err := dao.UpdateCollection(userId, videoId, utils.IsCollection); err != nil {
						log.Printf(err.Error())
						flag = true //出现问题
					}
				}
				//一遍流程下来正常执行了，那就打断结束，不再尝试
				if flag == false {
					break
				}
			}
		}
	}
}

// consumerCollectionDel 赞关系删除的消费方式。
func (l *CollectionMQ) consumerCollectionDel(messages <-chan amqp.Delivery) {
	for d := range messages {
		// 参数解析。
		params := strings.Split(fmt.Sprintf("%s", d.Body), " ")
		userId := params[0]
		videoId, _ := strconv.ParseInt(params[1], 10, 64)
		//最多尝试操作数据库的次数
		for i := 0; i < utils.Attempts; i++ {
			flag := false //默认无问题
			//取消赞行为，只有当前状态是点赞状态才会发起取消赞行为，所以如果查询到，必然是cancel==0(点赞)
			//先查询是否有这条数据
			collectionInfo, err := dao.GetCollectionInfo(userId, videoId)
			//如果有问题，说明查询数据库失败，返回错误信息err:"get collectionInfo failed"
			if err != nil {
				log.Printf(err.Error())
				flag = true //出现问题
			} else {
				if collectionInfo == (dao.Collection{}) { //只有当前是点赞状态才能取消点赞这个行为
					// 所以如果查询不到数据则返回错误信息:"can't find data,this action invalid"
					log.Printf(errors.New("can't find data,this action invalid").Error())
				} else {
					//如果查询到数据，则更新为取消赞状态
					//如果有问题，说明插入数据库失败，打印错误信息err:"update data fail"
					if err := dao.UpdateCollection(userId, videoId, utils.UnCollection); err != nil {
						log.Printf(err.Error())
						flag = true
					}
				}
			}
			//一遍流程下来正常执行了，那就打断结束，不再尝试
			if flag == false {
				break
			}
		}
	}
}

var RmqCollectionAdd *CollectionMQ
var RmqCollectionDel *CollectionMQ

// InitCollectionRabbitMQ 初始化rabbitMQ连接。
func InitCollectionRabbitMQ() {
	RmqCollectionAdd = NewCollectionRabbitMQ("collection_add")
	go RmqCollectionAdd.Consumer()

	RmqCollectionDel = NewCollectionRabbitMQ("collection_del")
	go RmqCollectionDel.Consumer()
}

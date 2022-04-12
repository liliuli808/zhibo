package kafka

import (
	"context"
	"errors"
	"github.com/Shopify/sarama"
	"log"
	"time"
	"zhibo/utils"
)

type Config struct {
	Topic   string `yaml:"topic"`
	Address string `yaml:"address"`
}

const CompressionCode = 3
const MaxMessageBytes = 51200000

// Product 生产者
type Product struct {
	Config Config
	Client sarama.SyncProducer
}

func (k *Product) Instance() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll              // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner // 环形选择分区，在所有分区中循环选择一个
	config.Producer.Return.Successes = true                       // 成功交付的消息将在success channel返回
	config.Producer.Compression = CompressionCode                 // 压缩 lz4 兼容线上0.1版本
	config.Producer.MaxMessageBytes = MaxMessageBytes             // 单消息体限制大小50m

	// 连接 kafka
	client, err := sarama.NewSyncProducer([]string{k.Config.Address}, config)
	if err != nil {
		panic(err)
	}
	k.Client = client
}

func (k *Product) Push(topic string, value string) error {
	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(value)

	_, _, err := k.Client.SendMessage(msg)
	if err != nil {
		return err
	}
	return nil
}

// bulk 批量发送
func (k *Product) bulk(topic string, messageBulk map[string]string) error {
	var bulk []*sarama.ProducerMessage
	for key, message := range messageBulk {
		ProducerMessage := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(key),
			Value: sarama.StringEncoder(message),
		}
		bulk = append(bulk, ProducerMessage)
	}

	err := k.Client.SendMessages(bulk)
	if err != nil {
		return err
	}
	return nil
}

func (k *Product) Close() error {
	return k.Client.Close()
}

// ConsumerGroupHandler 消费组Handler
type ConsumerGroupHandler struct {
	Handler func(m *sarama.ConsumerMessage, count int)
}

func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {
		// 执行handler定义
		h.Handler(msg, len(claim.Messages()))
		// 若无异常，则确认消费
		sess.MarkMessage(msg, "")
		msg = nil
	}
	return nil
}

// ConsumerGroup 消费组
type ConsumerGroup struct {
	Config       Config
	GroupID      string
	Topics       []string
	Group        sarama.ConsumerGroup
	GroupHandler ConsumerGroupHandler
	Ready        chan bool
}

func (kcg *ConsumerGroup) initConsumerGroup() (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0 // specify appropriate version
	config.Consumer.Return.Errors = true
	config.Consumer.Fetch.Max = MaxMessageBytes
	config.Producer.Compression = CompressionCode
	config.Producer.MaxMessageBytes = MaxMessageBytes
	return sarama.NewConsumerGroup([]string{kcg.Config.Address}, kcg.GroupID, config)
}

// 重试
func (kcg *ConsumerGroup) reConnect() error {
	for i := 0; i < 10; i++ {
		group, err := kcg.initConsumerGroup()
		if err == nil {
			kcg.Group = group
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return errors.New("ReConnect Failed")
}

func (kcg *ConsumerGroup) Start(ctx context.Context) {

	group, err := kcg.initConsumerGroup()
	utils.PanicNotNil(err)
	kcg.Group = group
	go func() {
		for err := range kcg.Group.Errors() {
			log.Println("KafkaConsumerGroup ERROR", err)
			if err.Error() == "kafka: broker not connected" {
				err := kcg.reConnect()
				utils.PanicNotNil(err)
			}
			// 不中断
			//panic(err)
		}
	}()
	kcg.Group.Consume(ctx, kcg.Topics, kcg.GroupHandler)
	//panic(err)
}

func (kcg *ConsumerGroup) Close() {
	// 关闭消费组
	err := kcg.Group.Close()
	// 确认是否存在异常
	panic(err)
}

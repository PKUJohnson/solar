package toolkit

import (
	"fmt"

	"github.com/nsqio/go-nsq"
	std "github.com/PKUJohnson/solar/std"
)

// TODO consumer
func CreateNSQConsumer(topicName, channelName string, lookupAddressList []string,
	messageHandler nsq.Handler) *nsq.Consumer {
	// TODO
	configObj := nsq.NewConfig()
	consumerInstance, err := nsq.NewConsumer(topicName, channelName, configObj)
	if err != nil {
		content := fmt.Sprintf("CreateNSQConsumer NewConsumer err:%s topic:%s channel:%s",
			err, topicName, channelName)
		panic(content)
	}
	consumerInstance.AddHandler(messageHandler)
	consumerInstance.ChangeMaxInFlight(2000) // TODO https://github.com/nsqio/go-nsq/issues/179
	if err := consumerInstance.ConnectToNSQLookupds(lookupAddressList); err != nil {
		content := fmt.Sprintf("CreateNSQConsumer ConnectToNSQLookupds err:%s topic:%s channel:%s",
			err, topicName, channelName)
		panic(content)
	}
	return consumerInstance
}

func CreateConcurrentNSQConsumer(topicName, channelName string, lookupAddressList []string,
	messageHandler nsq.Handler, count int) []*nsq.Consumer {
	// TODO
	clientList := make([]*nsq.Consumer, 0)
	for i := 0; i < count; i += 1 {
		nsqClient := CreateNSQConsumer(topicName, channelName, lookupAddressList, messageHandler)
		clientList = append(clientList, nsqClient)
	}
	return clientList
}

func CloseNSQConsumer(client *nsq.Consumer) {
	client.Stop()
}

func CloseNSQConsumerList(clientList []*nsq.Consumer) {
	for _, client := range clientList {
		CloseNSQConsumer(client)
	}
}

// TODO produce

func CreateProducer(nsqdAddress string) *nsq.Producer {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(nsqdAddress, config)
	if err != nil {
		fmt.Println("FBI WARNING BUG BUG CreateProducer NewProducer Error...", nsqdAddress, err, producer)
	} else {
		producer.Ping()
	}
	return producer
}

func CreateConsumer(cnf *std.ConfigNsqConsumer) (*nsq.Consumer, error) {
	c := nsq.NewConfig()
	return nsq.NewConsumer(cnf.TopicName, cnf.ChannelName, c)
}

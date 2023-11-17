package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

func main() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Consumer.Return.Errors = true

	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalln(err)
		log.Fatalln("11111111")
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalln(err)
		log.Fatalln("222222")
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition("to-ender", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalln(err)
		log.Fatalln("333333")
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		for {
			message := &sarama.ProducerMessage{
				Topic: "to-ender",
				Value: sarama.StringEncoder("Hello, World!"),
			}

			partition, offset, err := producer.SendMessage(message)
			if err != nil {
				log.Printf("Failed to send message: %v\n", err)
			} else {
				fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)
			}

			time.Sleep(1 * time.Second)
		}
	}()

	// 消费者获取消息并打印
	for msg := range partitionConsumer.Messages() {
		fmt.Printf("Message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)
		fmt.Printf("Message value:%s\n", string(msg.Value))
	}
}

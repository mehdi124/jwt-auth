package email

import (
	"fmt"
	"context"
	"log"
	"github.com/go-redis/redis/v8"
	"github.com/keithwachira/go-taskq"
	"encoding/json"
)


var streamName = "send_emails"




//func NewEmailReceivedFromClient(data interface{}) {
//
//	rdq := Redis.NewClient()
//	err := rdq.XAdd(context.Background(), &redis.XAddArgs{
//		///this is the name we want to give to our stream
//		///in our case we called it send_order_emails
//		//note you can have as many stream as possible
//		//such as one for email...another for notifications
//		Stream:       streamName,
//		MaxLen:       0,
//		MaxLenApprox: 0,
//		ID:           "",
//		//values is the data you want to send to the stream
//		//in our case we send a map with email and message keys
//		Values: data,
//	}).Err()
//	if err != nil {
//		log.Fatal("failed ")
//		return
//	}
//	fmt.Fprintf(w, `We received you order`)
//}


type RedisStreamsProcessing struct {
	Redis *redis.Client
	//other dependencies e.g. logger database goes here
}

func (r *RedisStreamsProcessing) Process(job interface{}) {

	//the go redis client returns the redis stream data as type [redis.XMessage]
	if data, ok := job.(redis.XMessage); ok {

		emailTemplate := data.Values["template"]

		message := data.Values["value"]
		messageValue , _ :=  message.(string)
		var emailData EmailData
		json.Unmarshal([]byte( messageValue ),&emailData)

		fmt.Printf("I am sending an email to the email  %v and type is  %T  %v \n ", data,emailData.Email,emailTemplate)
		//here we can decide to delete each entry when it is processed
		//in that case you can use the redis xdel command i.e:

		SendEmail(&emailData,emailTemplate.(string))

		//TODO handle job failed and success jobs for repeat or delete
		r.Redis.XDel(context.Background(),streamName,data.ID).Err()


	} else {
		log.Println("wrong type of data sent")
	}
}

func StartProcessingEmails(rdb *redis.Client) {
	redisStreams := RedisStreamsProcessing{
		Redis: rdb,
	}
	//in this case we have started 5 goroutines so at any moment we will
	//be sending a maximum of 5 emails.
	//you can adjust these parameters to increase or reduce
	q := taskq.NewQueue(5, 10, redisStreams.Process)
	//call startWorkers  in a different goroutine otherwise it will block
	go q.StartWorkers()
	//with our workers running now we can start listening to new events from redis stream
	//we start from id 0 i.e. the first item in the stream
	id := "0"
	for {

		var ctx = context.Background()
		data, err := rdb.XRead(ctx, &redis.XReadArgs{
			Streams: []string{streamName, id},
			//count is number of entries we want to read from redis
			Count: 4,
			//we use the block command to make sure if no entry is found we wait
			//until an entry is found
			Block: 0,
		}).Result()
		if err != nil {
			log.Println(err)
			log.Fatal(err)
		}
		///we have received the data we should loop it and queue the messages
		//so that our jobs can start processing
		for _, result := range data {
			for _, message := range result.Messages {
				///we use EnqueueJobBlocking to send out jobs to the workers
				q.EnqueueJobBlocking(message)
				//here we set a new start id because we don't want to process old emails
				//so we have set the id to the last id we saw
				id = message.ID
			}
		}
	}
}
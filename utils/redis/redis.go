package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
	"strconv"
)

var ctx = context.Background()

func StoreVerificationCode(code string,user_id uint){

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	id := strconv.Itoa(user_id)
	err := rdb.Set(ctx, "verification_user_id_"+id, code, 10 * time.Second).Err()
	if err != nil {
		panic(err)
	}

}

func CheckVerificationCode(user_id uint,code string){

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	id := strconv.Itoa(user_id)
	val, err := rdb.Get(ctx, "verification_user_id_"+id).Result()
	if err != nil {
		panic(err)
	}

	if val != code{
		panic("invalid verification code")
	}

}

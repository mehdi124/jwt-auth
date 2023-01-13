package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var ctx = context.Background()

func StoreVerificationCode(user_id,code string){

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "verification_user_id_"+user_id, code, 10 * time.Second).Err()
	if err != nil {
		panic(err)
	}

}

func CheckVerificationCode(user_id ,code string){

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	val, err := rdb.Get(ctx, "verification_user_id_"+user_id).Result()
	if err != nil {
		panic(err)
	}

	if val != code{
		panic("invalid verification code")
	}

}

package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"task/configs"
	"task/db"
	"task/internal/flood_control"
)

func main() {
	cnf, err := configs.InitConfig("configs.yaml")
	if err != nil {
		log.Fatal(err)
	}

	addr := cnf.Connect.Host + ":" + cnf.Connect.Port
	opt := redis.Options{
		Addr:     addr,
		Password: cnf.Connect.Password,
		DB:       cnf.Connect.Database,
		Protocol: cnf.Connect.Protocol,
	}

	client, err := db.NewClient(&opt)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	fch := flood_control.NewFloodControlHandler(client, cnf.Interval, cnf.Limit)

	//failed flood control
	for i := 0; i < 12; i++ {
		client.AddNewMessage(ctx, 1, cnf.Interval)
	}
	control, _ := fch.Check(ctx, 1)
	if control == false {
		fmt.Printf("User %d failed\n", 1)
	}

	//pass flood control
	for i := 0; i < 3; i++ {
		client.AddNewMessage(ctx, 2, cnf.Interval)
	}
	control, _ = fch.Check(ctx, 2)
	if control == true {
		fmt.Printf("User %d pass\n", 2)
	}

}

// FloodControl интерфейс, который нужно реализовать.
// Рекомендуем создать директорию-пакет, в которой будет находиться реализация.
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}

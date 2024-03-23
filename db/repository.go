package db

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rClient *redis.Client
}

func NewClient(opt *redis.Options) (*Client, error) {
	rdb := redis.NewClient(opt)
	if rdb == nil {
		return nil, errors.New("can't create redis rClient")
	}
	return &Client{rClient: rdb}, nil
}
func (client *Client) GetAmountMessages1(ctx context.Context, userID uint64) (int, error) {
	val, err := client.rClient.Get(ctx, strconv.FormatUint(userID, 10)).Time()
	if err != nil {
		return -1, err
	}
	fmt.Println(val)
	return 1, nil
}

func (client *Client) GetAmountMessages(ctx context.Context, userID uint64) (int64, error) {
	key := fmt.Sprintf("messages:%d", userID)

	count, err := client.rClient.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return -1, err
	}

	return int64(len(count)), nil
}

func (client *Client) AddNewMessage(ctx context.Context, userID uint64, interval time.Duration) error {
	err := client.deleteMessages(ctx, userID, interval)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("messages:%d", userID)

	err = client.rClient.LPush(ctx, key, time.Now().Unix()).Err()
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) deleteMessages(ctx context.Context, userID uint64, interval time.Duration) error {
	key := fmt.Sprintf("messages:%d", userID)

	messages, err := client.rClient.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return err
	}

	var filteredMessages []string
	for _, message := range messages {
		timestamp, err := strconv.ParseInt(message, 10, 64)
		if err != nil {
			return err
		}
		if time.Now().Add(-interval).Unix() <= timestamp {
			filteredMessages = append(filteredMessages, message)
		}
	}

	err = client.rClient.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	for _, message := range filteredMessages {
		err = client.rClient.RPush(ctx, key, message).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

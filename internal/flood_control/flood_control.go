package flood_control

import (
	"context"
	"task/db"
	"time"
)

type FloodControlHandler struct {
	client   *db.Client
	interval time.Duration
	limit    int64
}

func NewFloodControlHandler(client *db.Client, interval time.Duration, limit int64) *FloodControlHandler {
	return &FloodControlHandler{
		client:   client,
		interval: interval,
		limit:    limit,
	}
}

func (fch *FloodControlHandler) Check(ctx context.Context, userID int64) (bool, error) {
	amount, err := fch.client.GetAmountMessages(ctx, userID)
	if err != nil {
		return false, err
	}
	if amount > fch.limit {
		return false, nil
	}
	return true, nil
}

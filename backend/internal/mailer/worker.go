package mailer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func StartEmailWorker(mailer *Mailer, rdb *redis.Client) {
	conn, err := mailer.dialer.Dial()
	if err != nil {
		logger.PrintFatal(fmt.Errorf("SMTP dial error: %w", err), nil)
	}

	defer conn.Close()
	for {
		res, err := rdb.BLPop(context.Background(), 0*time.Second, "email_queue").Result()
		if err != nil {
			logger.PrintError(fmt.Errorf("redis error: %w", err), nil)
			time.Sleep(2 * time.Second)
			continue
		}

		if len(res) < 2 {
			continue
		}

		var job EmailJob
		if err := json.Unmarshal([]byte(res[1]), &job); err != nil {
			logger.PrintError(fmt.Errorf("JSON decode error: %w", err), nil)
			continue
		}

		if err := mailer.SendUsingConn(conn, job); err != nil {
			logger.PrintError(fmt.Errorf("email worker error: %w", err), nil)
		}
	}
}

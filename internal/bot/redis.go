package bot

import (
	"errors"
)

const (
	NoExpirationTime = 0
)

var ErrKeyNotFound = errors.New("redis key not found")

// func (b *TGBot) SetUserState(userId int64, messageType messageId int64, page int) error {
// 	ctx := context.Background()
// 	// value := fmt.Sprintf("%d_%d", messageId, page)
//     status := b.Service.SetUserState(int(userId), )
// 	// status := b.redisClient.Set(ctx, strconv.FormatInt(userId, 10), value, NoExpirationTime)
// 	if status.Err() != nil {
// 		return fmt.Errorf("error while set user state: %w", status.Err())
// 	}
//
// 	return nil
// }

// func (b *TGBot) GetUserState(userId int64) (int, int, error) {
// 	ctx := context.Background()
// 	status := b.redisClient.Get(ctx, strconv.FormatInt(userId, 10))
//
// 	if err := status.Err(); err != nil {
// 		if err == redis.Nil {
// 			return 0, 0, ErrKeyNotFound
// 		}
// 		return 0, 0, fmt.Errorf("error while getting user state: %w", err)
// 	}
//
// 	result, err := status.Result()
// 	if err != nil {
// 		return 0, 0, fmt.Errorf("error while getting user state value: %w", err)
// 	}
//
// 	messageId, err := strconv.Atoi(strings.Split(result, "_")[0])
// 	if err != nil {
// 		return 0, 0, fmt.Errorf("error while parsing messageId from redis: %w", err)
// 	}
// 	page, err := strconv.Atoi(strings.Split(result, "_")[1])
// 	if err != nil {
// 		return 0, 0, fmt.Errorf("error while parsing page number from redis: %w", err)
// 	}
//
// 	return messageId, page, nil
// }

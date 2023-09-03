package rdb

import "context"

func (r *RedisStore) SaveMessage(ctx context.Context, message string, key string, set string) error {
	return r.DB.HSet(ctx, set, key, message).Err()
}
func (r *RedisStore) GetMessage(ctx context.Context, key string) (map[string]string, error) {
	return r.DB.HGetAll(ctx, key).Result()

}

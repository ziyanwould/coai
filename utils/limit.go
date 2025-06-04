package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// CheckIPLimit 检查某个IP在指定窗口内的操作次数限制
// scene: 业务场景（如"register"、"verify"），limit: 次数上限，window: 时间窗口
// 返回true表示未超限，可以继续，false表示超限
func CheckIPLimit(cache *redis.Client, ip string, scene string, limit int, window time.Duration) (bool, error) {
	ctx := context.Background()
	ipKey := fmt.Sprintf("%s:ip:%s", scene, ip)
	count, _ := cache.Get(ctx, ipKey).Int()
	if count >= limit {
		return false, nil
	}
	pipe := cache.TxPipeline()
	pipe.Incr(ctx, ipKey)
	pipe.Expire(ctx, ipKey, window)
	_, err := pipe.Exec(ctx)
	return true, err
}

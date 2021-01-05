package rolling_counter

import (
	"log"
	"time"
)

type Option func(*options)

type options struct {
	interval  time.Duration
	bucketNum int64
}

// defaultOptions 默认配置
// 每隔500毫秒统计一次
// 统计10秒内的数据
var defaultOptions = options{
	interval:  time.Millisecond * 500,
	bucketNum: 20,
}

func WithInterval(interval time.Duration) Option {
	return func(o *options) {
		if interval <= 0 {
			log.Printf("rolling_counter: warning wrong interval %d", interval)
			return
		}
		o.interval = interval
	}
}

func WithBucketNum(bucketNum int64) Option {
	return func(o *options) {
		if bucketNum <= 0 {
			log.Printf("rolling_counter: warning wrong bucketNum %d", bucketNum)
			return
		}
		o.bucketNum = bucketNum
	}
}

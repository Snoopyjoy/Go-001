package rolling_counter

import (
	"sync"
	"time"
)

type RollingCounter interface {
	Incr(v float64)
	Sum(now time.Time) float64
	Max(now time.Time) float64
	Avg(now time.Time) float64
}

type rollingCounter struct {
	interval    time.Duration    // 每个桶的时间
	bucketNum   int64            // 桶的数量
	buckets     []float64        // 配合cursor循环利用buckets列表
	cursor      int              // 桶起始索引
	startNanoTs int64            // 开始计数纳秒时间戳
	mutex       *sync.RWMutex    // 读写锁
	nowFunc     func() time.Time // 当前时间获取方法 方便测试
}

func NewRollingCounter(options ...Option) RollingCounter {
	return newRollingCounter(options...)
}

func newRollingCounter(options ...Option) *rollingCounter {
	opt := defaultOptions

	for _, o := range options {
		o(&opt)
	}

	return &rollingCounter{
		interval:  opt.interval,
		bucketNum: opt.bucketNum,
		buckets:   make([]float64, opt.bucketNum),
		mutex:     &sync.RWMutex{},
		nowFunc:   time.Now,
	}
}

// updateBuckets 更新并清空过期的bucket
func (c *rollingCounter) updateBuckets(now time.Time) {

	endIndex := c.getCurrentIndex(now)

	// 无重叠区间
	if endIndex < 0 {
		c.resetBuckets(now)
		return
	}

	startIndex := endIndex + 1 - int(c.bucketNum)
	// 无重叠区间
	if startIndex >= int(c.bucketNum) {
		c.resetBuckets(now)
		return
	}

	if endIndex >= int(c.bucketNum) {
		// 清理过期bucket
		for i := 0; i < startIndex; i++ {
			c.buckets[c.getRealIndex(i)] = 0
		}
		c.startNanoTs += c.interval.Nanoseconds() * int64(startIndex)
		c.cursor = c.getRealIndex(startIndex)
	}
}

func (c *rollingCounter) getRealIndex(i int) int {
	return (i + c.cursor) % int(c.bucketNum)
}

func (c *rollingCounter) resetBuckets(now time.Time) {
	c.cursor = 0
	for i := range c.buckets {
		c.buckets[i] = 0
	}
	c.startNanoTs = now.UnixNano()
}

// getCurrentIndex 获取当前时间在滑动窗口中的索引
func (c *rollingCounter) getCurrentIndex(now time.Time) int {
	nowTS := now.UnixNano()
	endIndex := (nowTS - c.startNanoTs) / int64(c.interval)
	return int(endIndex)
}

// getIndexRange 获取当前时间在滑动窗口中的索引范围
func (c *rollingCounter) getIndexRange(now time.Time) [2]int {
	endIndex := c.getCurrentIndex(now)
	startIndex := endIndex + 1 - int(c.bucketNum)
	if startIndex < 0 {
		startIndex = 0
	}
	if endIndex >= int(c.bucketNum) {
		endIndex = int(c.bucketNum) - 1
	}
	if endIndex < 0 {
		endIndex = -1
	}
	return [2]int{startIndex, endIndex}
}

func (c *rollingCounter) Incr(v float64) {
	if v == 0 {
		return
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	now := c.nowFunc()
	c.updateBuckets(now)
	c.buckets[c.getRealIndex(c.getCurrentIndex(now))] += v
}

func (c *rollingCounter) Sum(now time.Time) float64 {
	sum := float64(0)
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// 起止窗口索引
	windowRange := c.getIndexRange(now)

	for i := windowRange[0]; i <= windowRange[1]; i++ {
		sum += c.buckets[c.getRealIndex(i)]
	}

	return sum
}

func (c *rollingCounter) Max(now time.Time) float64 {
	max := float64(0)
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// 起止窗口索引
	windowRange := c.getIndexRange(now)

	for i := windowRange[0]; i <= windowRange[1]; i++ {
		idx := c.getRealIndex(i)
		if c.buckets[idx] > max {
			max = c.buckets[idx]
		}
	}

	return max
}

func (c *rollingCounter) Avg(now time.Time) float64 {
	return c.Sum(now) / float64(c.bucketNum)
}

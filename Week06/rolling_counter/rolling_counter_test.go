package rolling_counter

import (
	"testing"
	"time"
)

func TestCounter(t *testing.T) {

	now := time.Unix(1609845889, 0)
	nowFunc := func() time.Time {
		return now
	}
	// 3个桶 每个桶间隔500毫秒
	counter := newRollingCounter(WithBucketNum(3), WithInterval(time.Millisecond*500))
	counter.nowFunc = nowFunc
	counter.Incr(1)
	counter.Incr(1)
	counter.Incr(1)

	// 500 毫秒后
	now = time.Unix(1609845889, int64(time.Millisecond)*500)
	counter.Incr(1)
	counter.Incr(1)

	// 1秒后
	now = time.Unix(1609845889, int64(time.Millisecond)*1000)
	counter.Incr(1)

	if counter.Sum(now) != 6 {
		t.Errorf("expect sum 6 but got %f", counter.Sum(now))
	}

	if counter.Max(now) != 3 {
		t.Errorf("expect max 3 but got %f", counter.Max(now))
	}

	if counter.Avg(now) != 2 {
		t.Errorf("expect avg 2 but got %f", counter.Avg(now))
	}

	// 1.5秒后 第一个桶废弃
	now = time.Unix(1609845889, int64(time.Millisecond)*1500)
	if counter.Sum(now) != 3 {
		t.Errorf("expect sum 3 but got %f", counter.Sum(now))
	}

	if counter.Max(now) != 2 {
		t.Errorf("expect max 2 but got %f", counter.Max(now))
	}

	if counter.Avg(now) != 1 {
		t.Errorf("expect avg 1 but got %f", counter.Avg(now))
	}

	// 循环利用第一个桶
	counter.Incr(6)

	if counter.Sum(now) != 9 {
		t.Errorf("expect sum 9 but got %f", counter.Sum(now))
	}

	if counter.Max(now) != 6 {
		t.Errorf("expect max 6 but got %f", counter.Max(now))
	}

	if counter.Avg(now) != 3 {
		t.Errorf("expect avg 3 but got %f", counter.Avg(now))
	}

	// 三秒后所有桶都废弃 重新计数
	now = time.Unix(1609845889, int64(time.Millisecond)*3000)
	counter.Incr(3)
	if counter.Sum(now) != 3 {
		t.Errorf("expect sum 3 but got %f", counter.Sum(now))
	}

	if counter.Max(now) != 3 {
		t.Errorf("expect max 3 but got %f", counter.Max(now))
	}

	if counter.Avg(now) != 1 {
		t.Errorf("expect avg 1 but got %f", counter.Avg(now))
	}
}

func TestGetIndexRange(t *testing.T) {
	now := time.Unix(1609845889, 0)
	nowFunc := func() time.Time {
		return now
	}
	// 10个桶，间隔500ms
	counter := newRollingCounter(WithBucketNum(10), WithInterval(time.Millisecond*500))
	counter.nowFunc = nowFunc

	counter.Incr(1)

	// 当前
	v := counter.getIndexRange(now)
	if v != [2]int{0, 0} {
		t.Errorf("expect index range [0,0] but got [%d,%d]", v[0], v[1])
	}

	// 1秒前 已过期
	v = counter.getIndexRange(time.Unix(1609845888, 0))
	if v != [2]int{0, -1} {
		t.Errorf("expect index range [0,-1] but got [%d,%d]", v[0], v[1])
	}

	// 过了3秒后
	v = counter.getIndexRange(time.Unix(1609845889, int64(time.Millisecond)*500*6))
	if v != [2]int{0, 6} {
		t.Errorf("expect index range [0,6] but got [%d,%d]", v[0], v[1])
	}

	// 过了6秒后
	v = counter.getIndexRange(time.Unix(1609845889, int64(time.Millisecond)*500*12))
	if v != [2]int{3, 9} {
		t.Errorf("expect index range [3,9] but got [%d,%d]", v[0], v[1])
	}
}

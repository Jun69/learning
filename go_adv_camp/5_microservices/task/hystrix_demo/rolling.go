package rolling


import (
	"sync"
	"time"
)


type Number struct {
	Mutex   *sync.RWMutex
	Buckets map[int64]*numberBucket
	cnt  int64   // 保留总数
}

type numberBucket struct {
	Value float64
}


func NewNumber(num int64) *Number {
	r := &Number{
		Mutex:   &sync.RWMutex{},
		Buckets: make(map[int64]*numberBucket),
		cnt: num,
	}
	return r
}

func (r *Number) getCurrentBucket() *numberBucket {
	now := time.Now().Unix()
	var bucket *numberBucket
	var ok bool

	if bucket, ok = r.Buckets[now]; !ok {
		bucket = &numberBucket{}
		r.Buckets[now] = bucket
	}

	return bucket
}

func (r *Number) removeOldBuckets() {
	now := time.Now().Unix() - r.cnt

	for timestamp := range r.Buckets {
		// TODO: configurable rolling window
		if timestamp <= now {
			delete(r.Buckets, timestamp)
		}
	}
}

// 只能操作当前桶
func (r *Number) Increment(i float64) {
	if i == 0 {
		return
	}

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	b := r.getCurrentBucket()
	b.Value += i
	r.removeOldBuckets()
}

// UpdateMax updates the maximum value in the current bucket.
func (r *Number) UpdateMax(n float64) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	b := r.getCurrentBucket()
	if n > b.Value {
		b.Value = n
	}
	r.removeOldBuckets()
}

func (r *Number) Sum(now time.Time) float64 {
	sum := float64(0)

	r.Mutex.RLock()
	defer r.Mutex.RUnlock()
	for timestamp, bucket := range r.Buckets {
		if timestamp >= now.Unix()-r.cnt {
			sum += bucket.Value
		}
	}

	return sum
}

func (r *Number) Max(now time.Time) float64 {
	var max float64

	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	for timestamp, bucket := range r.Buckets {
		if timestamp >= now.Unix()-r.cnt {
			if bucket.Value > max {
				max = bucket.Value
			}
		}
	}

	return max
}

func (r *Number) Avg(now time.Time) float64 {
	return r.Sum(now) / float64(r.cnt)
}
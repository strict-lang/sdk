package concurrent

import (
	"sync"
	"sync/atomic"
)

type Latch interface {
	CountDown()
	Wait()
}

type latch struct {
	count int32
	lock *sync.Mutex
	condition *sync.Cond
}

func NewLatch(count int) Latch {
	lock := &sync.Mutex{}
	return &latch{
		count: int32(count),
		lock: lock,
		condition: sync.NewCond(lock),
	}
}

func (latch *latch) decrement() int32 {
	value := atomic.LoadInt32(&latch.count)
	target := value - 1
	for !atomic.CompareAndSwapInt32(&latch.count, value, target) {}
	return target
}

func (latch *latch) CountDown() {
	if latch.decrement() == 0 {
		latch.lock.Lock()
		defer latch.lock.Unlock()
		latch.condition.Broadcast()
	}
}

func (latch *latch) Wait() {
	latch.lock.Lock()
	defer latch.lock.Unlock()
	for atomic.LoadInt32(&latch.count) != 0 {
		latch.condition.Wait()
	}
}

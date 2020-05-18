package concurrent

import "testing"

func TestLatch(testing *testing.T) {
	latch := NewLatch(5)
	for count := 0; count < 5; count++ {
		go latch.CountDown()
	}
	// latch.Wait()
}

func TestLatch_WaitWhenZero(testing *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			testing.Error(err)
		}
	}()
	latch := NewLatch(5)
	for count := 0; count < 5; count++ {
		latch.CountDown()
	}
	// latch.Wait()
}

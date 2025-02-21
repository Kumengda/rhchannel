package mqtt

import (
	"sync"
)

type messageQueue struct {
	queue []ReportMessage
	mu    sync.Mutex
	cond  *sync.Cond
}

func newMessageQueue(bufferSize int) *messageQueue {
	mq := &messageQueue{}
	mq.cond = sync.NewCond(&mq.mu)
	return mq
}

func (mq *messageQueue) send(msg ReportMessage) {
	mq.mu.Lock()
	mq.queue = append(mq.queue, msg)
	mq.cond.Signal()
	mq.mu.Unlock()
}

func (mq *messageQueue) receive() interface{} {
	for {
		mq.mu.Lock()
		for len(mq.queue) == 0 {
			mq.cond.Wait()
		}
		msg := mq.queue[0]
		mq.queue = mq.queue[1:]
		mq.mu.Unlock()
		return msg
	}

}

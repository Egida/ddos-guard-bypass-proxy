package flaresolverr

import (
	"sync"
	"time"
)

type Lock struct {
	cb      func()
	mux     sync.Mutex
	timeout time.Duration
	timer   *time.Timer
}

func NewLock(timeout time.Duration, cb func()) *Lock {
	l := &Lock{
		timeout: timeout,
		cb:      cb,
	}

	return l
}

func (l *Lock) TryLock() bool {
	ok := l.mux.TryLock()
	if ok {
		l.cleanupTimeout()
	}

	return ok
}

func (l *Lock) UnLock() {
	l.mux.Unlock()
	l.startTimeout()
}

func (l *Lock) cleanupTimeout() {
	if l.timer == nil {
		return
	}

	l.timer.Stop()
}

func (l *Lock) startTimeout() {
	l.timer = time.AfterFunc(l.timeout, l.cb)
}

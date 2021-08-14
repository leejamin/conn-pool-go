package pool

import (
	"sync/atomic"
	"time"
)

type Conn interface {
	Close() error
}

type conn struct {
	Conn

	usedAt   int64
	pooled   bool
	createdAt time.Time
}

func NewConn(c Conn) *conn {
	cn := &conn{
		Conn:     c,
		createdAt: time.Now(),
	}
	cn.SetUsedAt(time.Now())
	return cn
}

func (cn *conn) UsedAt() time.Time {
	unix := atomic.LoadInt64(&cn.usedAt)
	return time.Unix(unix, 0)
}

func (cn *conn) SetUsedAt(tm time.Time) {
	atomic.StoreInt64(&cn.usedAt, tm.Unix())
}

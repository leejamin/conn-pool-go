package pool

import (
	"net"
	"sync/atomic"
	"time"
)

type Conn struct {
	net.Conn

	usedAt   int64
	pooled   bool
	createdAt time.Time
}

func NewConn(netConn net.Conn) *Conn {
	cn := &Conn{
		Conn:     netConn,
		createdAt: time.Now(),
	}
	cn.SetUsedAt(time.Now())
	return cn
}

func (cn *Conn) UsedAt() time.Time {
	unix := atomic.LoadInt64(&cn.usedAt)
	return time.Unix(unix, 0)
}

func (cn *Conn) SetUsedAt(tm time.Time) {
	atomic.StoreInt64(&cn.usedAt, tm.Unix())
}

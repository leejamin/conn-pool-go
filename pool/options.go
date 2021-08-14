package pool

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"runtime"
	"strings"
	"time"
)

type Options struct {
	// The network type, either tcp or unix.
	// Default is tcp.
	Network string
	// host:port address.
	Addr string

	// Dialer creates new network connection and has priority over
	// Network and Addr options.
	Dialer  func(context.Context) (Conn, error)

	OnClose func(Conn) error

	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout time.Duration
	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Use value -1 for no timeout and 0 for default.
	// Default is 3 seconds.
	ReadTimeout time.Duration
	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.
	// Default is ReadTimeout.
	WriteTimeout time.Duration

	PoolSize           int
	MinIdleConns       int
	MaxConnAge         time.Duration
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration

	// TLS Config to use. When set TLS will be negotiated.
	TLSConfig *tls.Config

	Logger *log.Logger
}

func (opt *Options) init() {
	if opt.Addr == "" {
		opt.Addr = "localhost"
	}

	if opt.Network == "" {
		if strings.HasPrefix(opt.Network, "/") {
			opt.Network = "unix"
		} else {
			opt.Network = "tcp"
		}
	}

	if opt.DialTimeout == 0 {
		opt.DialTimeout = 5 * time.Second
	}

	if opt.Dialer == nil {
		opt.Dialer = func(ctx context.Context) (Conn, error) {
			netDialer := &net.Dialer{
				Timeout: opt.DialTimeout,
				KeepAlive: 5 * time.Minute,
			}
			if opt.TLSConfig == nil {
				return netDialer.DialContext(ctx, opt.Network, opt.Addr)
			}
			return tls.DialWithDialer(netDialer, opt.Network, opt.Addr, opt.TLSConfig)
		}
	}

	if opt.PoolSize == 0 {
		opt.PoolSize = 10 * runtime.NumCPU()
	}
	switch opt.ReadTimeout {
	case -1:
		opt.ReadTimeout = 0
	case 0:
		opt.ReadTimeout = 3 * time.Second
	}
	switch opt.WriteTimeout {
	case -1:
		opt.WriteTimeout = 0
	case 0:
		opt.WriteTimeout = opt.ReadTimeout
	}
	if opt.PoolTimeout == 0 {
		opt.PoolTimeout = opt.ReadTimeout + time.Second
	}

	if opt.IdleTimeout == 0 {
		opt.IdleTimeout = 5 * time.Minute
	}
	if opt.IdleCheckFrequency == 0 {
		opt.IdleCheckFrequency = time.Minute
	}

	if opt.Logger == nil {
		opt.Logger = log.Default()
	}
}
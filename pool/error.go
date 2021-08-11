package pool

import (
	"context"
	"net"
)

func isBadConn(err error, allowTimeout bool) bool {
	switch err {
	case nil:
		return false
	case context.Canceled, context.DeadlineExceeded:
		return true
	}

	if allowTimeout {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return !netErr.Temporary()
		}
	}

	return true

}

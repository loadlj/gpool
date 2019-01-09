package pool

import (
	"sync"
	"errors"
)

var (
	ErrorClosed = errors.New("pool has closed")
	LimitNonPositive = errors.New("limit is nonPositive")
)

type IPool interface {
	Acquire() (IConn, error)
	Release(IConn)
	GetPayload() interface{}
	SetPayload(interface{})
	SetLimit(int)
	GetLimit() int
	Close()
}

type Pool struct {
	payload interface{}
	connPool chan IConn
	factory func(IPool) (IConn, error)
	exitChan chan struct{}
	limit int
	isClosed bool
	mutex sync.RWMutex
}

func (pool *Pool) GetLimit() int {
	return pool.limit
}

func (pool *Pool) SetLimit(limit int) {
	pool.limit = limit
}

func (pool *Pool) GetPayload() interface{} {
	return pool.payload
}

func (pool *Pool) SetPayload(payload interface{}) {
	pool.payload = payload
}

func (pool *Pool) Acquire() (conn IConn, err error) {
	pool.mutex.RLock()
	defer pool.mutex.RUnlock()
	if pool.isClosed {
		err = ErrorClosed
		return
	}
	conn = <- pool.connPool
	if conn == nil {
		conn, err = pool.factory(pool)
	}
	return
}

func (pool *Pool) Release(conn IConn) {
	// reject new conn when pool is closed
	pool.mutex.RLock()
	defer pool.mutex.RUnlock()
	if conn == nil {
		return
	}
	select {
	// pool is full
	case pool.connPool <- conn:
		return
	default:
		conn.Close()
	}
}

func (pool *Pool) Close() {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	conns := pool.connPool
	close(conns)
	for conn := range conns {
		if conn != nil {
			// try to close conn
			conn.Close()
		}
	}
	pool.isClosed = true
	return
}

func NewPool(limit int, factory func(IPool) (IConn, error)) (p IPool, err error) {
	if limit <= 0 {
		err = LimitNonPositive
		return
	}
	chanPool := make(chan IConn, limit)
	exitChan := make(chan struct{})
	// init chan pool
	for i := 0; i < limit; i ++ {
		chanPool <- nil
	}
	p = &Pool{
		limit: limit,
		connPool: chanPool,
		exitChan: exitChan,
		factory: factory,
	}
	return
}
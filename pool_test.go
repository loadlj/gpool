package pool

import (
	"testing"

	"github.com/go-redis/redis"
)

func newConnection(p IPool) (conn IConn, err error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	conn = NewConn(client, client.Close)
	return
}

func TestNewPool(t *testing.T) {
	_, err := NewPool(10, newConnection)
	if err != nil {
		t.Errorf("init pool err: %v", err)
	}
}


func TestPool_Acquire(t *testing.T) {
	p, _ := NewPool(1, newConnection)
	conn, err := p.Acquire()
	if err != nil {
		t.Errorf("Acquire conn err: %v", err)
		return
	}
	c := conn.GetClient().(*redis.Client)
	pong, err := c.Ping().Result()
	if err != nil {
		t.Errorf("redis client ping err: %v", err)
	}
	if pong != "PONG" {
		t.Errorf("redis client ping res is: %s", pong)
	}
	p.Release(conn)
}

func TestPool_Close(t *testing.T) {
	p, _ := NewPool(10, newConnection)
	p.Close()
	_, err := p.Acquire()
	if err != ErrorClosed {
		t.Errorf("close pool err %v", err)
	}
}


func BenchmarkPool_Acquire(b *testing.B) {
	b.ResetTimer()
	p, _ := NewPool(b.N, newConnection)
	for i := 0; i < b.N; i ++ {
		conn, _ := p.Acquire()
		p.Release(conn)
	}
	p.Close()
}
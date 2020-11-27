# gpool
Golang thread safe connection pool

## Install

```bash
go get github.com/loadlj/gpool
```

## Example

```go
func newConnection(p pool.IPool) (conn pool.IConn, err error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	conn = pool.NewConn(client, client.Close)
	return
}

func main(){
	var wg sync.WaitGroup
	p, _ := pool.NewPool(10, newConnection)
	for i := 0; i < 1000; i ++ {
		wg.Add(1)
		go func() {
			if conn, err := p.Acquire(); err == nil {
				c := conn.GetClient().(*redis.Client)
				c.IncrBy("test:key", 1)
				p.Release(conn)
				wg.Done()
			} else {
				fmt.Println(err)
			}
		}()
	}
	wg.Wait()
	p.Close()
}
```

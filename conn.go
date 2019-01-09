package pool

type IConn interface {
	GetClient() interface{}
	Close() error
}


type Connection struct {
	client interface{}
	close func() error
	isClosed bool
}

func (conn *Connection) GetClient() interface{} {
	return conn.client
}

func (conn *Connection) Close() error {
	if !conn.isClosed {
		err := conn.close()
		if err != nil {
			conn.isClosed = true
		}
		return err
	}
	return nil
}

func NewConn(client interface{}, close func() error) *Connection {
	return &Connection{client:client, close: close}
}
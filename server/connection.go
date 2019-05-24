package server

import (
	"fmt"
	"github.com/RobinUS2/tsxdb/rpc"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

func (instance *Instance) StartListening() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", instance.opts.ListenPort))
	if err != nil {
		return err
	}
	instance.SetRpcListener(listener)

	go func() {
		for {
			conn, err := listener.Accept()
			isShuttingDown := atomic.LoadInt32(&instance.shuttingDown) == 1
			if err != nil {
				if !isShuttingDown {
					log.Printf("%s", err)
				}
				break
			}

			go instance.ServeConn(conn)

			if isShuttingDown {
				break
			}
		}

		// close
		if err := instance.closeRpc(); err != nil {
			panic(err)
		}
	}()

	return nil
}

var closeMux sync.RWMutex

func (instance *Instance) closeRpc() error {
	closeMux.Lock()
	defer closeMux.Unlock()
	listener := instance.RpcListener()
	if listener != nil {
		if err := listener.Close(); err != nil {
			return err
		}
		instance.SetRpcListener(nil)
	}
	return nil
}

func (instance *Instance) ServeConn(conn net.Conn) {
	// register connection
	instance.RegisterConn(conn)
	atomic.AddInt64(&instance.pendingRequests, 1)

	// auth timeout
	go func() {
		time.Sleep(60 * time.Second)
		//log.Printf("check auth from %v", conn.RemoteAddr())
		_ = conn.Close()
		instance.RemoveConn(conn)
	}()

	// buffered writer
	srv := rpc.NewGobServerCodec(conn)

	// serve
	instance.rpc.ServeCodec(srv)

	// unregister
	atomic.AddInt64(&instance.pendingRequests, -1)
}

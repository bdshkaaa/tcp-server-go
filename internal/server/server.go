package server

import (
	"log"
	"net"
	"sync"

	"github.com/bdshkaaa/tcp-server-go/internal/worker"
)

type Server struct {
	addr     string
	listener net.Listener
	pool     *worker.Pool
	wg       sync.WaitGroup
	shutdown chan struct{}
}

func New(addr string, workerCount int) *Server {
	return &Server{
		addr:     addr,
		pool:     worker.NewPool(workerCount),
		shutdown: make(chan struct{}),
	}
}

func (s *Server) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.pool.Start()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.shutdown:
				return nil
			default:
				log.Printf("Accept error: %v", err)
				continue
			}
		}
		s.pool.Submit(conn)
	}
}

func (s *Server) Stop() {
	close(s.shutdown)
	s.listener.Close()
	s.pool.Stop()
	log.Println("Server gracefully stopped")
}

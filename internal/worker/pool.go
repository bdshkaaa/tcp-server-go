package worker

import (
	"log"
	"net"
	"sync"
)

type Pool struct {
	workers int
	tasks   chan net.Conn
	wg      sync.WaitGroup
	quit    chan struct{}
}

func NewPool(workers int) *Pool {
	return &Pool{
		workers: workers,
		tasks:   make(chan net.Conn),
		quit:    make(chan struct{}),
	}
}

func (p *Pool) Start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
	log.Printf("Worker pool started with %d workers", p.workers)
}

func (p *Pool) worker(id int) {
	defer p.wg.Done()
	for {
		select {
		case conn := <-p.tasks:
			p.handle(conn, id)
		case <-p.quit:
			log.Printf("Worker %d shutting down", id)
			return
		}
	}
}

func (p *Pool) handle(conn net.Conn, workerID int) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("Worker %d: read error: %v", workerID, err)
		return
	}

	msg := string(buf[:n])
	log.Printf("Worker %d received: %s", workerID, msg)

	response := "echo: " + msg + "\n"
	conn.Write([]byte(response))
}

func (p *Pool) Submit(conn net.Conn) {
	p.tasks <- conn
}

func (p *Pool) Stop() {
	close(p.quit)
	p.wg.Wait()
	close(p.tasks)
	log.Println("Worker pool stopped")
}

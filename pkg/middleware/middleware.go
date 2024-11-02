package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/m-ghobadi/BatchWise/pkg/config"
	"github.com/m-ghobadi/BatchWise/pkg/models"
)

type Middleware struct {
	Config            *config.Config
	EventInput        chan models.Event
	HighPriorityQueue chan models.Event
	LowPriorityQueue  chan models.Event
	BatchQueue        chan models.Batch
	StopChan          chan bool
	LastBatchTime     time.Time
	WorkerQueues      []chan models.Event // Added for Round Robin
	FirstRequestTime  time.Time
	LastRequestTime   time.Time
	RequestCount      int32
}

func NewMiddleware(cfg *config.Config) *Middleware {
	return &Middleware{
		Config:            cfg,
		EventInput:        make(chan models.Event, 1000),
		HighPriorityQueue: make(chan models.Event, 1000),
		LowPriorityQueue:  make(chan models.Event, 10000),
		BatchQueue:        make(chan models.Batch, 100),
		StopChan:          make(chan bool),
		LastBatchTime:     time.Now(),
	}
}

func (m *Middleware) Start() {
	var wg sync.WaitGroup
	wg.Add(4)

	// Start event handler
	go func() {
		defer wg.Done()
		m.eventHandler()
	}()

	// Start high-priority processor
	go func() {
		defer wg.Done()
		m.processHighPriorityEvents()
	}()

	// Start batch processor
	go func() {
		defer wg.Done()
		m.processBatchEvents()
	}()

	// Start processing interval optimizer
	go func() {
		defer wg.Done()
		m.processingIntervalHandler()
	}()

	// Start HTTP server
	http.HandleFunc("/event", m.httpHandler)
	http.HandleFunc("/metrics", m.reportHandler)
	fmt.Println("Starting BatchWise server on port 8181")
	http.ListenAndServe(":8181", nil)

	wg.Wait()
}

// // ---------------------FIFO Event Handler---------------------
func (m *Middleware) fifoEventHandler() {
	for {
		select {
		case event := <-m.EventInput:
			m.processEventFIFO(event)
		case <-m.StopChan:
			return
		}
	}
}

// Process event without prioritization or batching
func (m *Middleware) processEventFIFO(event models.Event) {
	// Simulate processing time
	go m.forwardHandler(event)
}

func (m *Middleware) StartFIFO() {
	var wg sync.WaitGroup
	wg.Add(1)

	// Start FIFO event handler
	go func() {
		defer wg.Done()
		m.fifoEventHandler()
	}()

	// Start HTTP server
	http.HandleFunc("/event", m.httpHandler)
	http.HandleFunc("/metrics", m.reportHandler)
	fmt.Println("Starting FIFO server on port 8181")
	http.ListenAndServe(":8181", nil)

	wg.Wait()
}

// // ----------------RR Event Handler---------------------
func (m *Middleware) startRoundRobinWorkers() {
	m.WorkerQueues = make([]chan models.Event, m.Config.WorkerCount)
	for i := 0; i < m.Config.WorkerCount; i++ {
		m.WorkerQueues[i] = make(chan models.Event, 1000)
		go m.roundRobinWorker(i)
	}
}

func (m *Middleware) roundRobinWorker(workerID int) {
	for event := range m.WorkerQueues[workerID] {
		m.processEventRoundRobin(event)
	}
}

func (m *Middleware) distributeEventsRoundRobin() {
	workerID := 0
	for {
		select {
		case event := <-m.EventInput:
			m.WorkerQueues[workerID] <- event
			workerID = (workerID + 1) % len(m.WorkerQueues)
		case <-m.StopChan:
			return
		}
	}
}

func (m *Middleware) processEventRoundRobin(event models.Event) {
	// Simulate processing time
	go m.forwardHandler(event)
}

func (m *Middleware) StartRoundRobin() {
	var wg sync.WaitGroup
	wg.Add(2)

	// Start Round Robin workers
	m.startRoundRobinWorkers()

	// Start event distributor
	go func() {
		defer wg.Done()
		m.distributeEventsRoundRobin()
	}()

	// Start HTTP server
	http.HandleFunc("/event", m.httpHandler)
	http.HandleFunc("/metrics", m.reportHandler)
	fmt.Println("Starting Round Robin server on port 8181")
	http.ListenAndServe(":8181", nil)

	wg.Wait()
}

// // ------------------Static Batch Event Handler---------------------
func (m *Middleware) staticBatchEventHandler() {
	var batch []models.Event
	for {
		select {
		case event := <-m.EventInput:
			batch = append(batch, event)
			if len(batch) >= m.Config.StaticBatchSize {
				m.ProcessStaticBatch(batch)
				batch = []models.Event{}
			}
		case <-m.StopChan:
			return
		}
	}
}

func (m *Middleware) ProcessStaticBatch(batch []models.Event) {
	// Process batch of events
	// Simulate batch processing time
	for _, event := range batch {
		go m.forwardHandler(event)
	}
}

func (m *Middleware) StartStaticBatch() {
	var wg sync.WaitGroup
	wg.Add(1)

	// Start static batch event handler
	go func() {
		defer wg.Done()
		m.staticBatchEventHandler()
	}()

	// Start HTTP server
	http.HandleFunc("/event", m.httpHandler)
	http.HandleFunc("/metrics", m.reportHandler)
	fmt.Println("Starting Static Batch server on port 8181")
	http.ListenAndServe(":8181", nil)

	wg.Wait()
}

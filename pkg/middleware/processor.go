package middleware

import (
	"github.com/m-ghobadi/BatchWise/pkg/models"
)

func (m *Middleware) processHighPriorityEvents() {
	for {
		select {
		case event := <-m.HighPriorityQueue:
			m.processRealTime(event)
		case <-m.StopChan:
			return
		}
	}
}

func (m *Middleware) processRealTime(event models.Event) {

	// Implement actual processing logic here
	// fmt.Printf("Processing high-priority event: %s\n", event.ID)
	// // Simulate processing time
	// time.Sleep(50 * time.Millisecond)
	// TODO: Add logic to forward the event to the appropriate microservice
	go m.forwardHandler(event)
}

func (m *Middleware) processBatch(batch models.Batch) {
	// Implement actual batch processing logic here
	// fmt.Printf("Processing batch of %d events\n", len(batch.Events))
	// // Simulate batch processing time
	// time.Sleep(time.Duration(len(batch.Events)) * 10 * time.Millisecond)
	// TODO: Add logic to forward the batch of events to the appropriate microservice
	for _, event := range batch.Events {
		go m.forwardHandler(event)
	}
}

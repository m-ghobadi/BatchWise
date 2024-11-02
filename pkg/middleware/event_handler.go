package middleware

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/m-ghobadi/BatchWise/pkg/models"
)

func (m *Middleware) httpHandler(w http.ResponseWriter, r *http.Request) {
	if m.FirstRequestTime.IsZero() {
		m.FirstRequestTime = time.Now()
	}
	m.LastRequestTime = time.Now()

	// Parse event urgency from request
	eventUrgency, _ := strconv.ParseFloat(r.Header.Get("X-Event-Urgency"), 64)

	// Parse event from request (simplified)
	event := models.Event{
		ID:           r.Header.Get("X-Event-ID"),
		Type:         r.Header.Get("X-Event-Type"),
		Urgency:      eventUrgency,
		ReceivedTime: time.Now(),
		// Payload: r.Body,
		Request: r,
	}
	if event.Urgency > 3 {
		event.IsUserHighPriority = true
	}

	// Add event to the input channel
	m.EventInput <- event

	w.WriteHeader(http.StatusOK)
}

func (m *Middleware) eventHandler() {
	for {
		select {
		case event := <-m.EventInput:
			// Duration limiting logic
			if m.LastRequestTime.Sub(m.FirstRequestTime) >= time.Duration(5)*time.Minute {
				log.Fatal("Total Request: ", m.RequestCount)
			}

			// System metrics (should be collected from monitoring tools)
			systemLoad := getSystemLoad()
			traffic := getTrafficPatterns()

			// Classify event
			m.classifyEvent(&event, systemLoad, traffic)

			// Add event to appropriate queue
			if event.Priority >= m.Config.Thresholds.Priority {
				event.IsSysteHighPriority = true
				m.HighPriorityQueue <- event
			} else {
				m.LowPriorityQueue <- event
			}
		case <-m.StopChan:
			return
		}
	}
}

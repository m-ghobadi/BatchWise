package middleware

import (
	"github.com/m-ghobadi/BatchWise/pkg/models"
)

func (m *Middleware) classifyEvent(event *models.Event, systemLoad, traffic float64) {
	w := m.Config.Weights
	event.Priority = w.W1*getTypeScore(event.Type) + w.W2*event.Urgency + w.W3*systemLoad + w.W4*traffic
}

func getTypeScore(eventType string) float64 {
	switch eventType {
	case "transaction":
		return 10.0
	case "command":
		return 9.0
	case "log":
		return 8.0
	case "notification":
		return 7.0
	case "query":
		return 6.0
	default:
		return 1.0
	}
}

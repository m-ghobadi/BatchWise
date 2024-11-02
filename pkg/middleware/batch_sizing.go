package middleware

import (
	"math"
	"time"

	"github.com/m-ghobadi/BatchWise/pkg/models"
)

func (m *Middleware) calculateOptimalBatchSize(P_avg, L, T float64) int {
	constants := m.Config.Constants
	B_opt := (constants.Alpha*P_avg + constants.Beta) / (L + constants.Gamma*T)
	B_opt = clamp(B_opt, m.Config.BatchSizeLimits.Min, m.Config.BatchSizeLimits.Max)
	return int(B_opt)
}

func (m *Middleware) batchSizingAlgorithm(events []models.Event, L, T float64) models.Batch {
	var P_sum float64
	for _, event := range events {
		P_sum += event.Priority
	}

	P_avg := P_sum / float64(len(events))
	B_opt := m.calculateOptimalBatchSize(P_avg, L, T)

	// Create batch
	batchSize := int(math.Min(float64(len(events)), float64(B_opt)))
	batchEvents := events[:batchSize]

	return models.Batch{
		Events: batchEvents,
	}
}

func (m *Middleware) processBatchEvents() {
	var lowPriorityEvents []models.Event
	ticker := time.NewTicker(time.Microsecond) // Check every Nanosecond
	for {
		select {
		case event := <-m.LowPriorityQueue:
			lowPriorityEvents = append(lowPriorityEvents, event)
			// Check if batch size reached
			if len(lowPriorityEvents) >= int(m.Config.BatchSizeLimits.Max) {
				m.processCurrentBatch(&lowPriorityEvents)
			}
		case <-ticker.C:
			// Process batch based on processing interval
			if len(lowPriorityEvents) > 0 {
				// Get the processing interval
				interval := m.getProcessingInterval()
				if time.Since(m.LastBatchTime) >= time.Duration(interval)*time.Microsecond {
					m.processCurrentBatch(&lowPriorityEvents)
				}
			}
		case <-m.StopChan:
			return
		}
	}
}

func (m *Middleware) processCurrentBatch(lowPriorityEvents *[]models.Event) {
	systemLoad := getSystemLoad()
	traffic := getTrafficPatterns()
	batch := m.batchSizingAlgorithm(*lowPriorityEvents, systemLoad, traffic)

	// Process the batch
	m.processBatch(batch)

	// Update last batch processing time
	m.LastBatchTime = time.Now()

	// Remove processed events from the queue
	if len(*lowPriorityEvents) > len(batch.Events) {
		*lowPriorityEvents = (*lowPriorityEvents)[len(batch.Events):]
	} else {
		*lowPriorityEvents = []models.Event{}
	}
}

package middleware

import (
	"sync/atomic"
	"time"
)

var processingInterval atomic.Value

func (m *Middleware) updateProcessingInterval(currentLoad, eventPriority float64) float64 {
	constants := m.Config.Constants
	I := constants.C / (currentLoad + eventPriority)
	I = clamp(I, m.Config.IntervalLimits.Min, m.Config.IntervalLimits.Max)
	processingInterval.Store(I)
	return I
}

func (m *Middleware) processingIntervalHandler() {
	ticker := time.NewTicker(time.Duration(m.Config.SamplingInterval) * time.Microsecond)
	for {
		select {
		case <-ticker.C:
			load := getSystemLoad()
			priority := getAveragePriority()
			m.updateProcessingInterval(load, priority)
		case <-m.StopChan:
			return
		}
	}
}

func (m *Middleware) getProcessingInterval() float64 {
	interval := processingInterval.Load()
	if interval == nil {
		return m.Config.IntervalLimits.Max
	}
	return interval.(float64)
}

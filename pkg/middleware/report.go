package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/m-ghobadi/BatchWise/pkg/metrics"
)

func (m *Middleware) reportHandler(w http.ResponseWriter, r *http.Request) {
	// Create a report
	eventLogs, totalCount := metrics.PrintTableOutput()

	fmt.Println("Total Request: ", totalCount, " Duration: ", m.LastRequestTime.Sub(m.FirstRequestTime))

	eventLogJSON, err := json.MarshalIndent(eventLogs, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the report
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(eventLogJSON)
}

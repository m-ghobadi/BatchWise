package middleware

import (
	"github.com/m-ghobadi/BatchWise/pkg/metrics"
	"github.com/m-ghobadi/BatchWise/pkg/models"
)

// Handler function to route requests
func (m *Middleware) forwardHandler(event models.Event) {

	metrics.LogEvent(event)

}

// // Service struct to define a service and its URL
// type Service struct {
// 	Name string
// 	URL  *url.URL
// }

// // Define your external services
// var services = map[string]Service{
// 	"transaction":  {Name: "Service transaction", URL: parseURL("http://localhost:5010")},
// 	"log":          {Name: "Service log", URL: parseURL("http://localhost:5020")},
// 	"notification": {Name: "Service notification", URL: parseURL("http://localhost:5030")},
// 	"command":      {Name: "Service command", URL: parseURL("http://localhost:5040")},
// 	"query":        {Name: "Service query", URL: parseURL("http://localhost:5030")},
// }

// // Function to parse URLs and handle errors
// func parseURL(rawURL string) *url.URL {
// 	u, err := url.Parse(rawURL)
// 	if err != nil {
// 		log.Fatalf("Error parsing URL: %s", err)
// 	}
// 	return u
// }
// // Handler function to route requests
// func (m *Middleware) forwardHandler(event models.Event) {

// 	serviceName := event.Type
// 	service, exists := services[serviceName]

// 	if !exists {
// 		log.Printf("Service %s not found", serviceName)
// 		return
// 	}

// 	// Create a new request to the target service
// 	targetURL := service.URL
// 	req, err := http.NewRequest(event.Request.Method, targetURL.String(), event.Request.Body)
// 	if err != nil {
// 		log.Printf("Error creating request: %s", err)
// 		return
// 	}

// 	// Copy headers from original request
// 	req.Header = event.Request.Header

// 	// Send the request to the external service
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Printf("Error forwarding request: %s", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	// Read response from external service
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Printf("Error reading response: %s", err)
// 		return
// 	}

// 	fmt.Printf("Forwarded event %s to %s, resp: %s, body: %s\n", event.ID, serviceName, resp.Status, string(body))

// 	if resp.StatusCode != http.StatusOK {
// 		log.Printf("Error response from %s: %s", serviceName, resp.Status)
// 		return

// 	}
// 	go metrics.SendEventForwardedNotification(event)

// 	// // Log the response
// 	// log.Printf("Response from %s: %s", serviceName, string(body))
// }

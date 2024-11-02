package metrics

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/m-ghobadi/BatchWise/pkg/models"
	"github.com/rodaine/table"
)

type EventType int

const (
	Transaction EventType = iota
	Log
	Notification
	Command
	Query
)

type EventTypeLog struct {
	ServiceType            string
	EventReceivedCount     int32
	HighPrioritySysteCount int32
	HighPriorityUserCount  int32
	HighPrioritySystemRate float64
	HighPriorityUserRate   float64
	TotalSHPHT             time.Duration
	TotalUHPHT             time.Duration
	TotalHoldingTime       time.Duration
	AverageHoldingTime     time.Duration
	AvgSHPHT               time.Duration
	AvgUHPHT               time.Duration
}

var EventTypeLogs = map[EventType]EventTypeLog{
	Transaction:  {},
	Log:          {},
	Notification: {},
	Command:      {},
	Query:        {},
}
var EventLogsList = []models.Event{}
var eventMux = sync.Mutex{}

type SystemMetrics struct {
	ReportTime time.Time
	CpuLoad    float64
	MemLoad    float64
}

type SystemMetricsLog struct {
	Duration   time.Duration
	AverageCPU float64
	AverageMem float64
}

var SystemMetricsLogs []SystemMetrics

//-------------Event Logs-------------------------------------------------------------------------------

func LogEvent(event models.Event) {
	// Log event
	eventMux.Lock()
	event.CompletedTime = time.Now()
	event.HoldingTime = event.CompletedTime.Sub(event.ReceivedTime)
	EventLogsList = append(EventLogsList, event)
	eventMux.Unlock()

}

func EventLogs() map[EventType]EventTypeLog {

	EventTypeLogs = map[EventType]EventTypeLog{
		Transaction:  {ServiceType: "Transaction"},
		Log:          {ServiceType: "Log"},
		Notification: {ServiceType: "Notification"},
		Command:      {ServiceType: "Command"},
		Query:        {ServiceType: "Query"},
	}

	for _, event := range EventLogsList {
		eventType := getStringTypeToEvent(event.Type)
		eventLog := EventTypeLogs[eventType]

		eventLog.EventReceivedCount++
		eventLog.TotalHoldingTime += event.HoldingTime

		if event.IsSysteHighPriority {
			eventLog.HighPrioritySysteCount++
			eventLog.TotalSHPHT += event.HoldingTime
		}
		if event.IsUserHighPriority {
			eventLog.HighPriorityUserCount++
			eventLog.TotalUHPHT += event.HoldingTime
		}

		EventTypeLogs[eventType] = eventLog
	}

	for eventType := range EventTypeLogs {
		eventLog := EventTypeLogs[eventType]
		if eventLog.EventReceivedCount > 0 {
			eventLog.HighPrioritySystemRate = float64(eventLog.HighPrioritySysteCount) / float64(eventLog.EventReceivedCount)
			eventLog.HighPriorityUserRate = float64(eventLog.HighPriorityUserCount) / float64(eventLog.EventReceivedCount)

			avgT := eventLog.TotalHoldingTime.Nanoseconds() / int64(eventLog.EventReceivedCount)
			if eventLog.HighPrioritySysteCount > 0 {
				avgHPSHT := eventLog.TotalSHPHT.Nanoseconds() / int64(eventLog.HighPrioritySysteCount)
				eventLog.AvgSHPHT = time.Duration(avgHPSHT)
			}

			if eventLog.HighPriorityUserCount > 0 {
				avgHPUHT := eventLog.TotalUHPHT.Nanoseconds() / int64(eventLog.HighPriorityUserCount)
				eventLog.AvgUHPHT = time.Duration(avgHPUHT)
			}
			eventLog.AverageHoldingTime = time.Duration(avgT)
			EventTypeLogs[eventType] = eventLog
		}
	}
	return EventTypeLogs
}
func getStringTypeToEvent(eventType string) EventType {
	switch eventType {
	case "transaction":
		return Transaction
	case "log":
		return Log
	case "notification":
		return Notification
	case "command":
		return Command
	case "query":
		return Query
	default:
		return -1
	}
}

func PrintTableOutput() (map[EventType]EventTypeLog, int32) {

	eventLogs := EventLogs()
	var totalRequest int32

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Event type", "Request count", "Sys P rate", "User P rate", "Avg. holding time", "Avg. Sys P", "Avg. User P")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, eventLog := range eventLogs {
		totalRequest += eventLog.EventReceivedCount
		tbl.AddRow(eventLog.ServiceType, eventLog.EventReceivedCount, eventLog.HighPrioritySystemRate, eventLog.HighPriorityUserRate,
			eventLog.AverageHoldingTime, eventLog.AvgSHPHT, eventLog.AvgUHPHT)
	}
	tbl.Print()

	return eventLogs, totalRequest

}

// -------------System Logs-------------------------------------------------------------------------------
func GetSystemMetrics() SystemMetrics {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	cpu, err := GetCPULoad()
	if err != nil {
		log.Printf("Error getting CPU load: %s", err)
	}

	return SystemMetrics{
		ReportTime: time.Now(),
		CpuLoad:    cpu,
		MemLoad:    float64(stats.Sys),
	}
}

func GetBasicSystemStats() {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	fmt.Printf("Number of CPUs: %d\n", runtime.NumCPU())
	fmt.Printf("Number of Goroutines: %d\n", runtime.NumGoroutine())
	fmt.Printf("Total Allocated Memory: %v MB\n", stats.TotalAlloc/1024/1024)
	fmt.Printf("System Memory: %v MB\n", stats.Sys/1024/1024)
}

func GetCPULoad() (float64, error) {
	// Read the contents of /proc/stat
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, err
	}

	// Split the data into lines
	lines := strings.Split(string(data), "\n")
	if len(lines) < 1 {
		return 0, fmt.Errorf("no data found")
	}

	// Parse the first line (which is the overall CPU)
	fields := strings.Fields(lines[0])
	if len(fields) < 5 {
		return 0, fmt.Errorf("not enough fields in /proc/stat")
	}

	// Extract user, nice, system, idle, and iowait
	user, _ := strconv.ParseFloat(fields[1], 64)
	nice, _ := strconv.ParseFloat(fields[2], 64)
	system, _ := strconv.ParseFloat(fields[3], 64)
	idle, _ := strconv.ParseFloat(fields[4], 64)
	iowait, _ := strconv.ParseFloat(fields[5], 64)

	// Total CPU time calculation
	total := user + nice + system + idle + iowait
	busy := total - idle

	// Calculate load as a percentage
	load := (busy / total) * 100.0

	return load, nil
}

func LogSystemMetrics() {
	for {
		SystemMetricsLogs = append(SystemMetricsLogs, GetSystemMetrics())
		time.Sleep(5 * time.Second)
	}
}

func GetSystemMetricsLogs() SystemMetricsLog {
	var totalCPUload, totalMemLoad float64
	for _, systemMetrics := range SystemMetricsLogs {
		// Calculate average CPU and Memory load
		totalCPUload += systemMetrics.CpuLoad
		totalMemLoad += systemMetrics.MemLoad
	}
	duration := SystemMetricsLogs[len(SystemMetricsLogs)-1].ReportTime.Sub(SystemMetricsLogs[0].ReportTime)
	averageCPU := totalCPUload / float64(len(SystemMetricsLogs))
	averageMem := totalMemLoad / float64(len(SystemMetricsLogs))

	return SystemMetricsLog{
		Duration:   duration,
		AverageCPU: averageCPU,
		AverageMem: averageMem,
	}
}

//--------------------------------------------------------------------------------------------

func SendEventProcessedNotification(event models.Event) {
	// Send notification
	req, err := http.NewRequest("GET", "http://127.0.0.1:8051/event_processed", nil)
	if err != nil {
		log.Printf("Error creating request: %s", err)
		return
	}
	req.Header.Set("X-Event-ID", event.ID)
	req.Header.Set("X-Event-Type", event.Type)
	req.Header.Set("X-Processed-Time", time.Now().Format(time.RFC3339))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending Event Processed notification: %s", err)
		return
	}
	defer resp.Body.Close()
}

func SendEventForwardedNotification(event models.Event) {
	// Send notification
	req, err := http.NewRequest("GET", "http://127.0.0.1:8051/event_forwarded", nil)
	if err != nil {
		log.Printf("Error creating request: %s", err)
		return
	}
	req.Header.Set("X-Event-ID", event.ID)
	req.Header.Set("X-Event-Type", event.Type)
	req.Header.Set("X-Forwarded-Time", time.Now().Format(time.RFC3339))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending Event Forwarded notification: %s", err)
		return
	}
	defer resp.Body.Close()
}

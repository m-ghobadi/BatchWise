package main

import (
	"flag"
	"fmt"

	"github.com/m-ghobadi/BatchWise/pkg/config"
	"github.com/m-ghobadi/BatchWise/pkg/middleware"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig("configs/config.yaml")

	// Initialize middleware
	mw := middleware.NewMiddleware(cfg)

	method := flag.String("m", "hybrid", "Method to start (hybrid, fifo, rr, static)")
	flag.Parse()

	switch *method {
	case "hybrid":
		mw.Start()
	case "fifo":
		mw.StartFIFO()
	case "rr":
		mw.StartRoundRobin()
	case "static":
		mw.StartStaticBatch()
	default:
		fmt.Println("Unknown method")
	}
}

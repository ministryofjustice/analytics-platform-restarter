package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

const defaultPort = 8000

var (
	logger *log.Logger
	port   int
)

func init() {
	logger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.LUTC|log.Lshortfile)

	port = getPort()
}

func main() {
	fmt.Printf("🔌 Starting server on port :%d...\n", port)
	server := NewServer(port)
	err := server.ListenAndServe()
	if err != nil {
		logger.Fatalf("🔻 Server shut down: %s", err)
	}
}

// getPort reads the HTTP server port from the PORT environment
// variable
func getPort() int {
	portStr, ok := os.LookupEnv("PORT")
	if !ok {
		logger.Printf("$PORT not set. Defaulting to %d.", defaultPort)
		return defaultPort
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		logger.Fatalf("$PORT must be an integer >= 1024: %s", err)
	}

	if port < 1024 {
		logger.Fatal("$PORT must be >= 1024.")
	}

	return port
}

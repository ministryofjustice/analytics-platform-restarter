package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	k8s "k8s.io/client-go/kubernetes"
)

const defaultPort = 8000

var (
	logger    *log.Logger
	port      int
	home      string
	k8sClient k8s.Interface
)

func init() {
	logger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.LUTC|log.Lshortfile)

	port = getPort()
	home = getHome()

	k8sClient = KubernetesClient(filepath.Join(home, ".kube", "config"))
}

func main() {
	log.Printf("⚡️ Starting server on port :%d...\n", port)
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
		logger.Fatalf("💥 $PORT must be an integer >= 1024: %s", err)
	}

	if port < 1024 {
		logger.Fatal("💥 $PORT must be >= 1024.")
	}

	return port
}

func getHome() string {
	home, ok := os.LookupEnv("HOME")
	if !ok {
		logger.Fatalf("💥 $HOME not set. It couldn't determine HOME directory.")
	}

	return home
}

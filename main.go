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
	namespace string
	home      string
	k8sClient k8s.Interface
)

func init() {
	logger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.LUTC|log.Lshortfile)

	readConfig()

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

func readConfig() {
	port = getPort()

	namespace = os.Getenv("NAMESPACE")
	if namespace == "" {
		logger.Fatalf("💥 $NAMESPACE not set or blank.")
	}

	home = os.Getenv("HOME")
	if home == "" {
		logger.Fatalf("💥 $HOME not set or blank.")
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

package go_playground

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// Task represents an event
type Task struct {
	ID   int
	Data string
}

// Worker function that listens for tasks
func worker(id int, jobs <-chan Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range jobs {
		fmt.Printf("Worker %d processing task %d with data: %s\n", id, task.ID, task.Data)
		time.Sleep(time.Second) // Simulating work
	}
}

// Global variables
var (
	jobQueue = make(chan Task, 100) // Buffered channel for tasks
	taskID   = 0
	mu       sync.Mutex
)

func eventHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	taskID++
	task := Task{ID: taskID, Data: "Event received"}
	mu.Unlock()

	jobQueue <- task // Send task to worker pool
	fmt.Fprintf(w, "Event %d added to queue\n", task.ID)
}

func main() {
	numWorkers := 5 // Fixed number of workers
	var wg sync.WaitGroup

	// Start worker pool
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobQueue, &wg)
	}

	// Set up HTTP server
	router := mux.NewRouter()
	router.HandleFunc("/event", eventHandler).Methods("POST")

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))

	wg.Wait() // Wait for workers (never actually reached since server runs forever)
}

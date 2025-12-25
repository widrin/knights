package actor

// Dispatcher manages the execution of actor message processing
type Dispatcher interface {
	// Schedule schedules a function to be executed
	Schedule(fn func())

	// Shutdown gracefully shuts down the dispatcher
	Shutdown()
}

// defaultDispatcher uses a fixed pool of goroutines
type defaultDispatcher struct {
	workers   int
	taskQueue chan func()
}

// NewDefaultDispatcher creates a new default dispatcher
func NewDefaultDispatcher(workers int) Dispatcher {
	d := &defaultDispatcher{
		workers:   workers,
		taskQueue: make(chan func(), 1000),
	}

	// Start worker goroutines
	for i := 0; i < workers; i++ {
		go d.worker()
	}

	return d
}

func (d *defaultDispatcher) worker() {
	for task := range d.taskQueue {
		task()
	}
}

func (d *defaultDispatcher) Schedule(fn func()) {
	select {
	case d.taskQueue <- fn:
	default:
		// Queue is full, execute in new goroutine
		go fn()
	}
}

func (d *defaultDispatcher) Shutdown() {
	close(d.taskQueue)
}

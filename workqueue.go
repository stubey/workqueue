// Package workqueue provides an implementation of a queue for performing
// tasks with a number of background worker processes. At its core, this
// package utilizes a lot of the inherent properties of channels.
package workqueue

// Work is a zero-argument function sendable to a WorkQueue
// The workqueue will execute the function (a callback)
type Work func()

// WorkQueue is a channel type that you can send a zero-argument Work function to.
type WorkQueue chan Work

// New creates a Workqueue (chan Work) and a dispatcher and returns a channel that you send a zero-argument function to.
// The dispatcher will make workers and listens on the returned channel for work requests and forwards them to a worker.
func New(ctx context.Context, numWorkers int) WorkQueue {
	queue := make(WorkQueue)
	d := make(dispatcher, numWorkers)
	go d.dispatch(queue)
	return queue
}

// dispatcher is a channel that receives worker channels that are ready to process
type dispatcher chan chan Work

// newDispatcher creates a Work channel of depth number of workers
// The New() fcn duplicates this functionality, or can create standalone dispatcher for te
func newDispatcher(queue WorkQueue, numWorkers int) dispatcher {
	d := make(dispatcher, numWorkers)
	go d.dispatch(queue)
	return d
}

// dispatch creates workers
// It then reads the input work queue for a task
// When it receives a task, it pulls a ready worker from the worker channel/queue and passes it the task
// If the input work queue is closed, it closes all of its workers input channels
func (d dispatcher) dispatch(queue WorkQueue) {
	// Create and start all of our workers.
	for i := 0; i < cap(d); i++ {
		w := make(worker)
		go w.work(d)
	}

	// Start the main loop in a goroutine.
	go func() {
		for work := range queue {
			go func(work Work) {
				worker := <-d
				worker <- work
			}(work)
		}

		// If we get here, the work queue has been closed, and we should
		// stop all of the workers.
		for i := 0; i < cap(d); i++ {
			w := <-d
			close(w)
		}
	}()
}

type worker chan Work

// work() passes itself to the dispatcher to let the dispatcher know it is ready for work
// it then waits in a goroutine for the dispatcher
func (w worker) work(d dispatcher) {
	// Add ourselves to the dispatcher.
	d <- w

	// Start the main loop.
	go w.wait(d)
}

// wait represents a ready goroutine
// it receives a work function from the dispatcher and runs it
// After running, it adds itself back to the dispatcher
func (w worker) wait(d dispatcher) {
	for work := range w {
		// Do the work.
		if work == nil {
			panic("nil work received")
		}

		work()

		// Re-add ourselves to the dispatcher.
		d <- w
	}
}

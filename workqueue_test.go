package workqueue_test

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stubey/workqueue"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestWorkQueue(t *testing.T) {
	var wg sync.WaitGroup

	for i := 10; i < 100; i *= 10 {
		wg.Add(i)
		q := workqueue.New(i)
		for j := 0; j < i; j++ {
			go func(w int) {
				q <- func() {
					dur := time.Duration(rand.Intn(10))
					time.Sleep(dur * time.Millisecond)
					t.Log(w)
					wg.Done()
				}
			}(j)
		}

		wg.Wait()
		close(q)
	}
}

func ExampleNew() {
	log.SetFlags(log.Lshortfile)
	// Create a new WorkQueue.
	wq := workqueue.New(1024)

	// This sync.WaitGroup is to make sure we wait until all of our work
	// is done.
	var wg sync.WaitGroup

	// Do some work.
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(v int) {
			wq <- func() {
				defer wg.Done()

				time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
				log.Println(v)
			}
		}(i)
	}

	// Wait for all of the work to finish, then close the WorkQueue.
	wg.Wait()
	close(wq)

	// Examples need Output to run, place at end of examples
	fmt.Println("ExampleNew Complete")
	// Output: ExampleNew Complete
}

func ExampleMe() {
	log.SetFlags(log.Lshortfile)

	numWorkers := 10

	ctx := context.Background()
	ctx, cancel = context.WithCancel(ctx)

	// Create a new WorkQueue.
	wq := workqueue.New(ctx, numWorkers)

	// This sync.WaitGroup is to make sure we wait until all of our work is done.
	var wg sync.WaitGroup

	// Do some work.
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(v int) {
			wq <- func() {
				defer wg.Done()

				time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
				log.Println(v)
			}
		}(i)
	}

	// Wait for all of the work to finish, then close the WorkQueue.
	wg.Wait()
	close(wq)

	// Examples need Output to run, place at end of examples
	fmt.Println("ExampleNew Complete")
	// Output: ExampleNew Complete
}

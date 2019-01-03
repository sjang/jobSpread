// Example provided with help from Jason Waldrip.
// Package work manages a pool of goroutines to perform work.
package work

import (
  "sync"
)

// Worker must be implemented by types that want to use
// the work pool.
type Worker interface {
  InitTask()
  DoTask()
  ErrorTask()
  PostTask()
}

// Pool provides a pool of goroutines that can execute any Worker
// tasks that are submitted.
type Pool struct {
  job chan Worker
  error chan string
  job_wg   sync.WaitGroup
  error_wg   sync.WaitGroup
}

var workerPool chan *Pool

// New creates a new work pool.
func New(maxGoroutines int) *Pool {
  p := Pool{
    job: make(chan Worker),
    error: make(chan string),
  }

  p.job_wg.Add(maxGoroutines)
  for i := 0; i < maxGoroutines; i++ {
    go func() {
      for w := range p.job {
        w.InitTask()
        w.DoTask()
        w.ErrorTask()
        w.PostTask()
      }
      p.job_wg.Done()
    }()

  }

  return &p
}

// Run submits work to the pool.
func (p *Pool) Run(w Worker) {
  p.job <- w
}

// Shutdown waits for all the goroutines to shutdown.
func (p *Pool) Shutdown() {
  close(p.job)
  close(p.error)
  p.job_wg.Wait()
  p.error_wg.Wait()
}

func CreateWorkerPool(maxPoolSize int) {
  workerPool = make(chan *Pool, maxPoolSize)
  maxGor := 30 // goroutine count per job
  for i := 0; i < maxPoolSize; i++ {
    newJob := New(maxGor)
    workerPool <- newJob
  }
}

func MyWorker() *Pool {
  var p *Pool
  select {
    case p = <-workerPool:
    //case <-time.After(10 * time.Millisecond): // 미리 생성된 초기개수 이상으로 필요한 경우
    default:
      p = New(10)
  }
  return p
}

func ReturnWorker(p *Pool) {
  select {
    case workerPool <- p :
    //case <-time.After(10 * time.Millisecond): // 초기값 이상으로 생성된 worker의 처리 -> Shotdown()
    default:
      p.Shutdown()
  }
}


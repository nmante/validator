package validator

import (
	"errors"
	"sync"
)

const (
	MinWorkers = 2
	MaxWorkers = 20000
)

var (
	ErrInvalidJobSlice = errors.New("Must pass in slice of Job interfaces.")
)

type Pool interface {
	Run()
	Work()
}

type WorkerPool struct {
	jobs       []Job
	jobsChan   chan Job
	waitGroup  sync.WaitGroup
	numWorkers int
}

func NewWorkerPool(numWorkers int, jobs interface{}, options ...func(*WorkerPool) error) (*WorkerPool, error) {
	js, ok := jobs.([]Job)
	if !ok {
		return nil, ErrInvalidJobSlice
	}

	pool := &WorkerPool{
		jobs:     js,
		jobsChan: make(chan Job),
	}

	pool.setNumWorkers(numWorkers)

	for _, option := range options {
		err := option(pool)

		if err != nil {
			return nil, err
		}
	}

	return pool, nil
}

func (p *WorkerPool) setNumWorkers(n int) {
	if n < MinWorkers {
		p.numWorkers = MinWorkers
	} else if n > MaxWorkers {
		p.numWorkers = MaxWorkers
	}

	p.numWorkers = n
}

func (p *WorkerPool) Work() {
	for job := range p.jobsChan {
		job.Run(&p.waitGroup)
	}
}

func (p *WorkerPool) AddJobs(jobs ...Job) {
	p.jobs = append(p.jobs, jobs...)
}

func (p *WorkerPool) Run() {
	for i := 0; i < p.numWorkers; i++ {
		go p.Work()
	}

	p.waitGroup.Add(len(p.jobs))

	for _, job := range p.jobs {
		p.jobsChan <- job
	}

	close(p.jobsChan)
	p.waitGroup.Wait()
}

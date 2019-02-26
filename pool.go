package validator

import (
	"sync"
)

const (
	MinWorkers = 2
	MaxWorkers = 20000
)

type Pool interface {
	Run()
	Work()
}

type workerPool struct {
	jobs       []Job
	jobsChan   chan Job
	waitGroup  sync.WaitGroup
	numWorkers int
}

func NewWorkerPool(numWorkers int, jobs []Job, options ...func(*workerPool) error) (*workerPool, error) {
	if jobs == nil {
		jobs = []Job{}
	}

	pool := &workerPool{
		jobs:     jobs,
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

func (p *workerPool) setNumWorkers(n int) {
	if n < MinWorkers {
		p.numWorkers = MinWorkers
	} else if n > MaxWorkers {
		p.numWorkers = MaxWorkers
	} else {
		p.numWorkers = n
	}
}

func (p *workerPool) Work() {
	for job := range p.jobsChan {
		job.Run(&p.waitGroup)
	}
}

func (p *workerPool) AddJobs(jobs ...Job) {
	p.jobs = append(p.jobs, jobs...)
}

func (p *workerPool) Run() {
	defer close(p.jobsChan)
	for i := 0; i < p.numWorkers; i++ {
		go p.Work()
	}

	p.waitGroup.Add(len(p.jobs))

	for _, job := range p.jobs {
		p.jobsChan <- job
	}

	p.waitGroup.Wait()
}

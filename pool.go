package main

import (
	"sync"
)

const (
	MinWorkers = 1
	MaxWorkers = 20000
)

type WorkerPool interface {
	Run()
	Work()
}

type FuncWorkerPool struct {
	jobs       []*FuncJob
	jobsChan   chan *FuncJob
	waitGroup  sync.WaitGroup
	numWorkers int
}

func NewFuncWorkerPool(numWorkers int, jobs []*FuncJob, options ...func(*FuncWorkerPool)) *FuncWorkerPool {
	n := numWorkers
	if numWorkers < MinWorkers {
		n = MinWorkers
	} else if numWorkers > MaxWorkers {
		n = MaxWorkers
	}

	pool := &FuncWorkerPool{
		numWorkers: n,
		jobs:       jobs,
		jobsChan:   make(chan *FuncJob),
	}

	for _, option := range options {
		option(pool)
	}

	return pool
}

func (p *FuncWorkerPool) work() {
	for job := range p.jobsChan {
		job.Run(&p.waitGroup)
	}
}

func (p *FuncWorkerPool) AddJob(job *FuncJob) {
	p.jobs = append(p.jobs, job)
}

func (p *FuncWorkerPool) Run() {
	for i := 0; i < p.numWorkers; i++ {
		go p.work()
	}

	p.waitGroup.Add(len(p.jobs))

	for _, job := range p.jobs {
		p.jobsChan <- job
	}

	close(p.jobsChan)
	p.waitGroup.Wait()
}

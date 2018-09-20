package main

import (
	"errors"
	"sync"
)

type Job interface {
	Run(wg *sync.WaitGroup)
}

type Result struct {
	isValid          bool
	validationErrors []string
	err              error
}

type FuncJob struct {
	value         interface{}
	validatorFunc Func
	Err           error
	Result        FuncResponse
}

func NewFuncJob(value interface{}, validatorFunc Func, options ...func(*FuncJob) error) (*FuncJob, error) {
	if value == nil {
		return &FuncJob{}, errors.New("Must pass a valid value")
	}

	funcJob := &FuncJob{
		value:         value,
		validatorFunc: validatorFunc,
	}

	for _, option := range options {
		err := option(funcJob)
		if err != nil {
			return &FuncJob{}, err
		}
	}

	return funcJob, nil
}

func (j *FuncJob) Run(wg *sync.WaitGroup) {
	fr, err := j.validatorFunc(j.value)
	j.Err = err
	j.Result = fr
	wg.Done()
}

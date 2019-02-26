package validator

import (
	"github.com/nmante/validator/funcs"
)

// Rule is a custom object that contains a key and validator functions
type Rule struct {
	Funcs          []funcs.Func
	Key            string
	IsRequired     bool
	EnableParallel bool
}

// RuleResponse is the result returned from executing all of the Funcs in a Rule. It includes
// useful information like the validation errors messages/strings, an error if any Funcs failed
// at runtime, and a boolean representing if the rule is valid
type RuleResponse struct {
	Key              string
	IsValid          bool
	ValidationErrors []string
}

func (r Rule) createFuncJobs(value interface{}) ([]Job, error) {
	jobs := []Job{}
	for _, f := range r.Funcs {
		job, err := NewFuncJob(value, f)
		if err != nil {
			return jobs, err
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (r Rule) execute(value interface{}) (RuleResponse, error) {
	errors := []string{}
	isValid := true

	jobs, err := r.createFuncJobs(value)
	if err != nil {
		return RuleResponse{}, err
	}

	if r.EnableParallel {
		pool, err := NewWorkerPool(len(jobs), jobs)
		if err != nil {
			return RuleResponse{}, err
		}
		pool.Run()
	}

	for _, job := range jobs {
		j, ok := job.(*FuncJob)
		if !ok {
			return RuleResponse{}, ErrMustBeFuncJob
		}

		if !r.EnableParallel {
			response, err := j.validatorFunc(value)
			j.Err = err
			j.Result = response
		}

		if j.Err != nil {
			return RuleResponse{}, j.Err
		}

		errors = append(errors, j.Result.Error)
		isValid = isValid && j.Result.IsValid
	}

	return RuleResponse{
		Key:              r.Key,
		ValidationErrors: errors,
		IsValid:          isValid,
	}, nil
}

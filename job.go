package validator

import (
	"errors"
	"sync"
)

var (
	ErrMustBeFuncJob = errors.New("Job must be of type FuncJob")
	ErrMustBeRuleJob = errors.New("Job must be of type RuleJob")
)

type Job interface {
	Run(wg *sync.WaitGroup)
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
	response, err := j.validatorFunc(j.value)
	j.Err = err
	j.Result = response
	wg.Done()
}

type RuleJob struct {
	value  interface{}
	rule   Rule
	Err    error
	Result RuleResponse
}

func NewRuleJob(value interface{}, rule Rule, options ...func(*RuleJob) error) (*RuleJob, error) {
	ruleJob := &RuleJob{
		value: value,
		rule:  rule,
	}

	for _, option := range options {
		err := option(ruleJob)
		if err != nil {
			return &RuleJob{}, err
		}
	}

	return ruleJob, nil
}

func (j *RuleJob) Run(wg *sync.WaitGroup) {
	response := j.rule.execute(j.value)
	j.Result = response
	wg.Done()
}

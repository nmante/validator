package main

// Rule is a custom object that contains a key and validator functions
type Rule struct {
	Funcs          []Func
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
	Err              error
}

func (r Rule) routine(value interface{}, responses chan<- RuleResponse) {
	responses <- r.execute(value)
}

func (r Rule) createJobs(value interface{}) ([]*FuncJob, error) {
	jobs := []*FuncJob{}
	for _, f := range r.Funcs {
		job, err := NewFuncJob(value, f)
		if err != nil {
			return jobs, err
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (r Rule) execute(value interface{}) RuleResponse {
	errors := []string{}
	isValid := true

	if r.EnableParallel {
		jobs, err := r.createJobs(value)
		if err != nil {
			return RuleResponse{
				Err: err,
			}
		}

		pool := NewFuncWorkerPool(len(jobs), jobs)
		pool.Run()

		for _, j := range jobs {
			if j.Err != nil {
				return RuleResponse{
					Key: r.Key,
					Err: j.Err,
				}
			}

			errors = append(errors, j.Result.Error)
			isValid = j.Result.IsValid
		}
	} else {
		for _, validatorFunc := range r.Funcs {
			response, err := validatorFunc(value)

			if err != nil {
				return RuleResponse{Err: err}
			}

			errors = append(errors, response.Error)
			isValid = response.IsValid
		}
	}

	return RuleResponse{
		Key:              r.Key,
		ValidationErrors: errors,
		IsValid:          isValid,
		Err:              nil,
	}
}

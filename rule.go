package validator

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

func (r Rule) createJobs(value interface{}) ([]Job, error) {
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

func (r Rule) execute(value interface{}) RuleResponse {
	errors := []string{}
	isValid := true

	jobs, err := r.createJobs(value)
	if err != nil {
		return RuleResponse{
			Err: err,
		}
	}

	if r.EnableParallel {
		pool, err := NewWorkerPool(len(jobs), jobs)
		if err != nil {
			return RuleResponse{
				Key: r.Key,
				Err: err,
			}
		}
		pool.Run()
	}

	for _, job := range jobs {
		j, ok := job.(*FuncJob)
		if !ok {
			return RuleResponse{
				Key: r.Key,
				Err: ErrMustBeFuncJob,
			}
		}

		if !r.EnableParallel {
			response, err := j.validatorFunc(value)
			if err != nil {
				j.Err = err
			} else {
				j.Result = response
			}
		}

		if j.Err != nil {
			return RuleResponse{
				Key: r.Key,
				Err: j.Err,
			}
		}

		errors = append(errors, j.Result.Error)
		isValid = isValid && j.Result.IsValid
	}

	return RuleResponse{
		Key:              r.Key,
		ValidationErrors: errors,
		IsValid:          isValid,
		Err:              nil,
	}
}

package main

// Func is function type that all validator functions must follow
type Func func(interface{}) (FuncResponse, error)

// FuncResponse contains info around if a validator function was valid. If it isn't valid, an
// error message is also returned
type FuncResponse struct {
	IsValid bool
	Error   string
}

// Response contains a bool for if all rules are valid, as well as error messages for invalid rules
type Response struct {
	Errors  map[string][]string
	IsValid bool
}

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

func (r Rule) execute(value interface{}) RuleResponse {
	errors := []string{}
	isValid := true

	if r.EnableParallel {
		type Job struct {
			v     interface{}
			vFunc Func
		}

		type Result struct {
			isValid          bool
			validationErrors []string
			err              error
		}

		jobs := make(chan Job, len(r.Funcs))
		results := make(chan Result, len(r.Funcs))

		worker := func(js <-chan Job, rs chan<- Result) {
			wErrors := []string{}
			wIsValid := true

			for job := range js {
				response, err := job.vFunc(job.v)

				if err != nil {
					results <- Result{err: err}
					return
				}

				wErrors = append(wErrors, response.Error)
				wIsValid = response.IsValid
			}

			rs <- Result{validationErrors: wErrors, isValid: wIsValid}
		}

		for i := 0; i < 3; i++ {
			go worker(jobs, results)
		}

		for _, f := range r.Funcs {
			jobs <- Job{v: value, vFunc: f}
		}

		close(jobs)

		for _ = range r.Funcs {
			result := <-results
			if result.err != nil {
				return RuleResponse{
					Key: r.Key,
					Err: result.err,
				}
			}

			errors = append(errors, result.validationErrors...)
			isValid = result.isValid
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

func (r Rule) routine(value interface{}, response chan<- RuleResponse) {
	response <- r.execute(value)
}

// Validator is an object that contains a set of rules that can be validated
type Validator struct {
	EnableParallel bool
	rules          map[string]Rule
}

// AddRule adds a rule to the validator
func (v Validator) AddRule(key string, funcs ...Func) {
	if rule, ok := v.rules[key]; ok {
		rule.Funcs = append(rule.Funcs, funcs...)
		return
	}

	v.rules[key] = Rule{Key: key, Funcs: funcs}
}

// Rules returns the list of rules for this validator object
func (v Validator) Rules() map[string]Rule {
	return v.rules
}

// New returns a validator object
func New(rules ...Rule) Validator {
	rs := map[string]Rule{}

	for _, rule := range rules {
		if _, ok := rs[rule.Key]; !ok {
			rs[rule.Key] = rule
			continue
		}

		r := rs[rule.Key]
		r.Funcs = append(r.Funcs, rule.Funcs...)
	}

	return Validator{
		rules: rs,
	}
}

// Validate runs all the rules of validation
func (v Validator) Validate(values map[string]interface{}) (Response, error) {
	errors := map[string][]string{}
	isValid := true

	for key, rule := range v.Rules() {
		if _, ok := values[key]; rule.IsRequired && !ok {
			errors[key] = []string{"is required"}
			isValid = false
		}
	}

	if !v.EnableParallel {
		for k, val := range values {
			if rule, ok := v.rules[k]; ok {
				cResponse := rule.execute(val)
				if cResponse.Err != nil {
					return Response{}, cResponse.Err
				}

				if !cResponse.IsValid {
					errors[cResponse.Key] = append(errors[cResponse.Key], cResponse.ValidationErrors...)
					isValid = false
				}
			}
		}
	} else {
		numValues := len(values)
		parallelResponses := make(chan RuleResponse, numValues)

		for k, val := range values {
			if rule, ok := v.rules[k]; ok {
				go rule.routine(val, parallelResponses)
			}
		}

		for i := 0; i < numValues; i++ {
			cResponse := <-parallelResponses
			if cResponse.Err != nil {
				return Response{}, cResponse.Err
			}

			if !cResponse.IsValid {
				errors[cResponse.Key] = append(errors[cResponse.Key], cResponse.ValidationErrors...)
				isValid = false
			}
		}
	}

	return Response{Errors: errors, IsValid: isValid}, nil
}

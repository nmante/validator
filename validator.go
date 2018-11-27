package validator

// Validator is an object that contains a set of rules that can be validated in parallel, or synchronously
type Validator struct {
	enableParallel bool
	rules          map[string]Rule
}

// Response contains a bool for if all rules are valid, as well as error messages for invalid rules
type Response struct {
	Errors  map[string][]string
	IsValid bool
}

// New returns a validator object
func New(rules []Rule, options ...Option) (*Validator, error) {
	rs := map[string]Rule{}

	for _, rule := range rules {
		if _, ok := rs[rule.Key]; !ok {
			rs[rule.Key] = rule
			continue
		}

		r := rs[rule.Key]
		r.Funcs = append(r.Funcs, rule.Funcs...)
		rs[rule.Key] = r
	}

	v := &Validator{
		enableParallel: false,
		rules:          rs,
	}

	for _, option := range options {
		err := option(v)

		if err != nil {
			return nil, err
		}
	}

	return v, nil
}

// AddRule adds a rule to the validator
func (v *Validator) AddRule(key string, funcs ...Func) *Validator {
	if rule, ok := v.rules[key]; ok {
		rule.Funcs = append(rule.Funcs, funcs...)
		v.rules[key] = rule
		return v
	}

	v.rules[key] = Rule{Key: key, Funcs: funcs}
	return v
}

// Rules returns the map of rules for this validator object
func (v *Validator) Rules() map[string]Rule {
	return v.rules
}

// Validate runs all the rules of validation
func (v *Validator) Validate(values map[string]interface{}) (Response, error) {
	errors := map[string][]string{}
	isValid := true
	jobs := []Job{}

	for key, rule := range v.Rules() {
		if value, ok := values[key]; ok {
			rj, err := NewRuleJob(value, rule)
			if err != nil {
				return Response{}, err
			}

			jobs = append(jobs, rj)
		}

		if _, ok := values[key]; rule.IsRequired && !ok {
			errors[key] = []string{"is required"}
			isValid = false
		}
	}

	if v.enableParallel {
		pool, err := NewWorkerPool(len(values), jobs)
		if err != nil {
			return Response{}, err
		}
		pool.Run()
	}

	for _, job := range jobs {
		j, ok := job.(*RuleJob)
		if !ok {
			return Response{}, ErrMustBeRuleJob
		}

		if !v.enableParallel {
			response, err := j.rule.execute(j.value)
			j.Err = err
			j.Result = response
		}

		if j.Err != nil {
			return Response{}, j.Err
		}

		if !j.Result.IsValid {
			errors[j.rule.Key] = append(errors[j.rule.Key], j.Result.ValidationErrors...)
			isValid = false
		}
	}

	return Response{Errors: errors, IsValid: isValid}, nil
}

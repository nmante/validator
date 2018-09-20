package main

// Validator is an object that contains a set of rules that can be validated in parallel, or synchronously
type Validator struct {
	EnableParallel bool
	rules          map[string]Rule
}

// Response contains a bool for if all rules are valid, as well as error messages for invalid rules
type Response struct {
	Errors  map[string][]string
	IsValid bool
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

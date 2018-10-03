package validator

import (
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	validator := New()
	if len(validator.Rules()) != 0 {
		t.Error("There should be 0 rules")
	}

	mustBeOne := func(v interface{}) (FuncResponse, error) {
		if v.(int) != 1 {
			return FuncResponse{false, "must be 1"}, nil
		}
		return FuncResponse{true, ""}, nil
	}
	v2 := New(Rule{Key: "random", Funcs: []Func{mustBeOne}})

	if len(v2.Rules()) != 1 {
		t.Error("There should be 1 rules")
	}

}

func TestValidate(t *testing.T) {
	validator := New()
	validator.AddRule("must_be_even", func(v interface{}) (FuncResponse, error) {
		val, ok := v.(int)
		if !ok {
			return FuncResponse{}, ErrTypeMismatch{v, "int"}
		}

		if val%2 != 0 {
			return FuncResponse{false, "must be even integer"}, nil
		}

		return FuncResponse{true, ""}, nil
	})

	response, _ := validator.Validate(map[string]interface{}{
		"must_be_even": 3,
	})

	if len(response.Errors) != 1 {
		t.Errorf("There should be 1 error, %+v", response)
	}
}

func TestRules(t *testing.T) {
	rules := []Rule{
		Rule{Key: "hello", Funcs: []Func{}},
		Rule{Key: "world", Funcs: []Func{}},
	}

	v := New(rules...)

	if len(rules) != len(v.Rules()) {
		t.Errorf("There should be %d rules", len(rules))
	}

	for _, r := range rules {
		rule, ok := v.Rules()[r.Key]
		if !ok {
			t.Errorf("Initialized keys don't match")
		} else if !reflect.DeepEqual(rule.Funcs, r.Funcs) {
			t.Errorf("Funcs don't match")
		}
	}
}

func TestAddRule(t *testing.T) {
	validator := New()
	validator.AddRule("must_be_odd", func(v interface{}) (FuncResponse, error) {
		val, ok := v.(int)
		if !ok {
			return FuncResponse{}, ErrTypeMismatch{v, "int"}
		}

		if val%2 == 0 {
			return FuncResponse{false, "must be an odd integer"}, nil
		}

		return FuncResponse{true, ""}, nil
	}).AddRule("must_be_odd", func(v interface{}) (FuncResponse, error) {
		val, ok := v.(int)
		if !ok {
			return FuncResponse{}, ErrTypeMismatch{v, "int"}
		}

		if val%2 == 0 {
			return FuncResponse{false, "must be an odd integer"}, nil
		}

		return FuncResponse{true, ""}, nil
	})

	if n := len(validator.Rules()); n != 1 {
		t.Errorf("There are %d rules. There should be 1.", n)
	}
	if n := len(validator.Rules()["must_be_odd"].Funcs); n != 2 {
		t.Errorf("There are %d funcs. There should be 2.", n)
	}
}

func TestIsStringInt(t *testing.T) {
	validator := New(
		Rule{
			Key: "page_size",
			Funcs: []Func{
				IsStringInt,
			},
		},
		Rule{
			Key: "username",
			Funcs: []Func{
				IsStringInt,
			},
		},
	)

	r, err := validator.Validate(map[string]interface{}{
		"page_size": "53",
		"username":  "abc",
	})

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(r.Errors) != 1 {
		t.Error("'username' value is not a string. There should be one validation error message")
	}
}

func TestIsStringBetweenInts(t *testing.T) {
	validator := New(
		Rule{
			Key: "page_size",
			Funcs: []Func{
				IsStringBetweenInts(10, 100),
			},
		},
		Rule{
			Key: "rando",
			Funcs: []Func{
				IsStringBetweenInts(1, 10),
			},
		},
	)

	r, err := validator.Validate(map[string]interface{}{
		"page_size": "53",
		"rando":     "5",
	})

	if err != nil {
		t.Errorf(err.Error())
	}

	if !r.IsValid {
		t.Error(r.Errors)
	}
}

func TestIsStringEqualToInt(t *testing.T) {
	validator := New(
		Rule{
			Key: "page_size",
			Funcs: []Func{
				IsStringEqualToInt(53),
			},
		},
	)

	r, err := validator.Validate(map[string]interface{}{
		"page_size": "54",
	})

	if err != nil {
		t.Error(err)
	}

	if r.IsValid {
		t.Error("IsStringEqualToInt requires matching values")
	}
}

// TestParallelPropertyValidation simulates long blocking calls. It should at the most take ~numSecondsBlock
// to complete
func TestParallelPropertyValidation(t *testing.T) {
	var numSecondsBlock time.Duration = 1
	type Asset struct {
		ID  int
		URL string
	}

	isVideoExists := func(v interface{}) (FuncResponse, error) {
		// Simulate long processing video processing task
		time.Sleep(numSecondsBlock * time.Second)
		return FuncResponse{true, ""}, nil
	}

	isImageExists := func(v interface{}) (FuncResponse, error) {
		// Simulate long processing image task
		time.Sleep(numSecondsBlock * time.Second)
		return FuncResponse{true, ""}, nil
	}

	// Call the isVideoExists function twice just to show both funcs process in parallel
	validator := New(
		Rule{
			Key:            "video",
			EnableParallel: true,
			Funcs: []Func{
				isVideoExists,
				isVideoExists,
			},
		},
		Rule{
			Key: "image",
			Funcs: []Func{
				isImageExists,
			},
		},
	)
	validator.EnableParallel = true

	video := Asset{
		ID:  1,
		URL: "https://s3.amazonaws.com/randomuri/hello.mp4",
	}

	image := Asset{
		ID:  1,
		URL: "https://s3.amazonaws.com/randomuri2/yo.jpg",
	}

	startTime := time.Now()
	r, err := validator.Validate(map[string]interface{}{
		"video": video,
		"image": image,
	})

	timeDuration := time.Since(startTime)

	if err != nil {
		t.Errorf(err.Error())
	}

	if !r.IsValid {
		t.Error(r.Errors)
	}

	min := numSecondsBlock * time.Second
	max := numSecondsBlock*time.Second + 100*time.Millisecond

	if !(min < timeDuration && timeDuration < max) {
		t.Errorf("Test should be between %v & %v. Actually took %v seconds", min, max, timeDuration)
	}
}

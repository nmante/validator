package validator

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/nmante/validator/compare"
	"github.com/nmante/validator/funcs"
	"github.com/nmante/validator/transform"
	"github.com/nmante/validator/types"
)

func TestNew(t *testing.T) {
	validator, _ := New(nil)
	if len(validator.Rules()) != 0 {
		t.Error("There should be 0 rules")
	}

	mustBeOne := func(v interface{}) (funcs.Response, error) {
		if v.(int) != 1 {
			return funcs.Response{false, "must be 1"}, nil
		}
		return funcs.Response{true, ""}, nil
	}

	v2, _ := New([]Rule{Rule{Key: "random", Funcs: []funcs.Func{mustBeOne}}})

	if len(v2.Rules()) != 1 {
		t.Error("There should be 1 rules")
	}

}

func TestValidate(t *testing.T) {
	validator, _ := New(nil)
	validator.AddRule("must_be_even", funcs.IsInt)
	validator.AddRule("must_be_even", func(v interface{}) (funcs.Response, error) {
		val, ok := v.(int)
		if !ok {
			return funcs.Response{}, errors.New("Not an integer")
		}

		if val%2 != 0 {
			return funcs.Response{false, "must be even integer"}, nil
		}

		return funcs.Response{true, ""}, nil
	})

	validator.AddRule("page_size", funcs.String.IsInt)
	validator.AddRule("num", funcs.IsTransformableTo(transform.StringToInt, types.Int))

	response, err := validator.Validate(map[string]interface{}{
		"must_be_even": 3,
		"page_size":    "100",
		"num":          "54",
	})

	if err != nil {
		t.Errorf("There should be no runtime errors, %+v", response)
	}

	if len(response.Errors) != 1 {
		t.Errorf("There should be 1 validation error, %+v", response)
	}
}

func TestRules(t *testing.T) {
	rules := []Rule{
		Rule{Key: "hello", Funcs: []funcs.Func{}},
		Rule{Key: "world", Funcs: []funcs.Func{}},
	}

	v, _ := New(rules)

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
	validator, _ := New(nil)
	validator.AddRule("must_be_odd", func(v interface{}) (funcs.Response, error) {
		val, ok := v.(int)
		if !ok {
			return funcs.Response{}, errors.New("not an integer")
		}

		if val%2 == 0 {
			return funcs.Response{false, "must be an odd integer"}, nil
		}

		return funcs.Response{true, ""}, nil
	}).AddRule("must_be_odd", func(v interface{}) (funcs.Response, error) {
		val, ok := v.(int)
		if !ok {
			return funcs.Response{}, errors.New("not an integer")
		}

		if val%2 == 0 {
			return funcs.Response{false, "must be an odd integer"}, nil
		}

		return funcs.Response{true, ""}, nil
	})

	if n := len(validator.Rules()); n != 1 {
		t.Errorf("There are %d rules. There should be 1.", n)
	}
	if n := len(validator.Rules()["must_be_odd"].Funcs); n != 2 {
		t.Errorf("There are %d funcs. There should be 2.", n)
	}
}

// TestParallelPropertyValidation simulates long blocking calls. It should at the most take ~numMillisecondsBlock
// to complete
func TestParallelPropertyValidation(t *testing.T) {
	var numMillisecondsBlock time.Duration = 250
	type Asset struct {
		ID  int
		URL string
	}

	isVideoExists := func(v interface{}) (funcs.Response, error) {
		// Simulate long processing video processing task
		time.Sleep(numMillisecondsBlock * time.Millisecond)
		return funcs.Response{true, ""}, nil
	}

	isImageExists := func(v interface{}) (funcs.Response, error) {
		// Simulate long processing image task
		time.Sleep(numMillisecondsBlock * time.Millisecond)
		return funcs.Response{true, ""}, nil
	}

	// Call the isVideoExists function twice just to show both funcs process in parallel
	validator, _ := New(
		[]Rule{
			Rule{
				Key:            "video",
				EnableParallel: true,
				Funcs: []funcs.Func{
					isVideoExists,
					isVideoExists,
				},
			},
			Rule{
				Key: "image",
				Funcs: []funcs.Func{
					isImageExists,
				},
			},
			Rule{
				Key: "page_size",
				Funcs: []funcs.Func{
					funcs.IsEqual(transform.StringToInt, compare.Int, 100),
				},
			},
		},
		OptionParallel(true),
	)

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
		"video":     video,
		"image":     image,
		"page_size": "100",
	})

	timeDuration := time.Since(startTime)

	if err != nil {
		t.Errorf(err.Error())
	}

	if !r.IsValid {
		t.Error(r.Errors)
	}

	min := numMillisecondsBlock * time.Millisecond
	max := (numMillisecondsBlock + 100) * time.Millisecond

	if !(min < timeDuration && timeDuration < max) {
		t.Errorf("Test should be between %v & %v. Actually took %v milliseconds", min, max, timeDuration)
	}
}

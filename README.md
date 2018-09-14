# Validator

## Overview

`validator` is a pluggable and parallelizable tool that can be used for validating any type of input.

## Usage

In a nutshell, you can first configure the validator with a list of M properties to validate. Each property can N validation `Func`s associated with it. Properties can be processed in parallel, AND the validation `Func`s associated with a property can also be processed in parallel.

The following sections gives a brief overview of how to use the validator.

### Initializing

You can either initialize a validator with a set of `Rule`s, or create an empty validator and add `Rule`s afterwards.

#### Initialize with rules

```go
import "validator"

...

paramValidator := validator.New(
	validator.Rule{
		Key: "page_size",
		IsRequired: true, // default is false
		Funcs: []validator.Func{
			validator.IsStringInt,
			validator.IsStringBetween(1, 100),
			func(v interface{}) (bool, string) {
				val := v.(string)
				if val == "50" {
					return true, ""
				}
				return false, "must equal 50"
			},
		},
	},
)
```

#### Add rules later

```go
import "validator"

...

queryValidator := validator.New()

...

queryValidator.AddRule("page_size", validator.IsStringInt)
queryValidator.AddRule("page_size", validator.IsStringBetween(1, 100))
queryValidator.AddRule("page_size", func(v interface{}) (validator.FuncResponse, error) {
	val, ok := v.(string)
	if !ok {
		return validator.FuncResponse{}, errors.New("my custom error message")
	}
	
	if val == "50" {
		return validator.FuncResponse{true, ""}, nil
	}
	
	return validator.FuncResponse{false, "must equal 50"}, nil
})
```

As you may have seen above, you can also pass in custom functions as long as they match this signature:

```go
func(v interface{}) (validator.FuncResponse, error)
```

### Validating your actual values

To validate your values, call the `.Validate` function (attached to the `validator` object with an argument of type `map[string]interface{}`:

```go
vr, err := paramValidator.Validate(map[string]interface{}{
	"page_size": "53",
})
```

This will return a `validator.Response` and an `error`. `validator.Response` contains a `bool` field which is `true` if validation was successful, or false otherwise. Any validation errors are returned in a `map[string][]string` where the keys represent the keys of the map you passed to `validator.Validate`. IF there was a any sort of unexpected `error` (i.e. not a validation error), this is returned as the second return argument.

Here's an example of how you might use the response from above:

```go
if err != nil {
	log.Println("Uh oh, unexpected error: ", err.Error())
	// return a 500 http response, exit program, or however you'd like to handle this
}
if !vr.IsValid {
	log.Println("Validation failed: ", vr.Errors)
	// return a 400 response, exit program, or however you'd like to handle this
}
```

### Parallel Validation

You can tell the validator to process your properties in parallel:

```go
v := validator.New(
	validator.Rule{
		Key: "video",
		Funcs: []validator.Func{
			longBlockingValidationFunc,
		},
	},
	validator.Rule{
		Key: "image",
		Funcs: []validator.Func{
			aLongBlockingImageValidationFunc,
		},
	},
)
v.EnableParallel = true // process each rule in parallel

...

vr, err := v.Validate(values)
```

You can also tell each `Rule` to process it's `Func`s in parallel

```go
import "validator"

...

paramValidator := validator.New(
	validator.Rule{
		Key: "video",
		IsRequired: false,
		EnableParallel: true, //Process the functions below in parallel
		Funcs: []validator.Func{
			longBlockingValidationFunc,
			anotherLongBlockingValidationFunc
		},
	},
)
```

package util

import "strings"

// AggregateError Aggregation error, merging multiple errors into one error
type AggregateError struct {
	errs []error
}

// NewAggregateError Create a new aggregate error structure
func NewAggregateError(errs []error) error {
	if nil == errs || len(errs) == 0 {
		return nil
	}
	return &AggregateError{
		errs: errs,
	}
}

func (e *AggregateError) Error() string {
	var temp []string
	for _, v := range e.errs {
		temp = append(temp, v.Error())
	}
	return strings.Join(temp, ", ")
}

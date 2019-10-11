/*
Copyright 2019 The Koala Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

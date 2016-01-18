package utils

import (
	"fmt"
	"reflect"
)

type FailedAssertion string

func (fa FailedAssertion) Error() string {
	return string(fa)
}

func Assert(truthy bool, msg string) {
	if !truthy {
		panic(FailedAssertion(msg))
	}
}

func AssertNil(value interface{}, msg string) {
	// NOTE:
	//   ValueOf returns the zero Value only if "value" is nil
	//   IsValid returns true unless called on the zero Value
	//
	if reflect.ValueOf(value).IsValid() {
		panic(FailedAssertion(fmt.Sprintf("%v (%v != nil)", msg, value)))
	}
}

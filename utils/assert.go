package utils

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
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

func InvalidCodePath() {
	pc, _, line, _ := runtime.Caller(1)
	path := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	caller := "Parser." + path[len(path)-1]
	message := fmt.Sprintf("Reached invalid code path at %v:%v\n", caller, line)
	panic(FailedAssertion(message))
}

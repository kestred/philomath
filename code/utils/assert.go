package utils

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
)

// Puts is useful for debugging to avoid adding and removing imports for fmt
// constantly if utils has already been imported
func Puts(value interface{}) {
	fmt.Printf("%v:%v\n", reflect.ValueOf(value).Type(), value)
}

func Errorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+msg+"\n", args...)
}

func Assert(truthy bool, msg string, args ...interface{}) {
	if !truthy {
		AssertionFailed(msg, args...)
	}
}

func AssertNil(value interface{}, msg string) {
	// NOTE:
	//   ValueOf returns the zero Value only if "value" is nil
	//   IsValid returns true unless called on the zero Value
	//
	if reflect.ValueOf(value).IsValid() {
		AssertionFailed("%v (%v != nil)", msg, value)
	}
}

func NotImplemented(what string) {
	AssertionFailed("\"%s\" hasn't been implemented", what)
}

func InvalidCodePath() {
	pc, _, line, _ := runtime.Caller(1)
	path := strings.Split(runtime.FuncForPC(pc).Name(), "/")
	caller := path[len(path)-1]
	AssertionFailed("Reached invalid code path at %v:%v", caller, line)
}

func AssertionFailed(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Assertion: "+msg+"\n", args...)

	// print call stack
	stackEnded := false
	fmt.Fprintf(os.Stderr, "\n")
	for i := 2; i < 30; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			stackEnded = true
			break
		}

		pathParts := strings.Split(file, "/")
		path := pathParts[len(pathParts)-1]
		if len(path) > 16 {
			path = "(*)" + path[len(path)-13:len(path)]
		}

		callerParts := strings.Split(runtime.FuncForPC(pc).Name(), "/")
		caller := callerParts[len(callerParts)-1]
		fmt.Fprintf(os.Stderr, "  %16s  %s:%d\n", path, caller, line)
	}
	fmt.Fprintf(os.Stderr, "\n")

	if !stackEnded {
		fmt.Fprintf(os.Stderr, "\nWhy things fail: your stack is to deep for sanity...\n")
	}
	os.Exit(1)
}

func Typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

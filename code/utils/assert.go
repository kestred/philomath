package utils

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
)

func Assert(truthy bool, msg string) {
	if !truthy {
		fmt.Println("Assertion:", msg)
		printStack()
		os.Exit(1)
	}
}

func AssertNil(value interface{}, msg string) {
	// NOTE:
	//   ValueOf returns the zero Value only if "value" is nil
	//   IsValid returns true unless called on the zero Value
	//
	if reflect.ValueOf(value).IsValid() {
		fmt.Printf("Assertion: %v (%v != nil)\n", msg, value)
		printStack()
		os.Exit(1)
	}
}

func InvalidCodePath() {
	pc, _, line, _ := runtime.Caller(1)
	path := strings.Split(runtime.FuncForPC(pc).Name(), "/")
	caller := path[len(path)-1]

	fmt.Printf("Assertion: Reached invalid code path at %v:%v\n", caller, line)
	printStack()
	os.Exit(1)
}

func printStack() {
	fmt.Println()
	for i := 2; i < 50; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			fmt.Println()
			return
		}

		pathParts := strings.Split(file, "/")
		path := pathParts[len(pathParts)-1]
		if len(path) > 16 {
			path = "(*)" + path[len(path)-14:len(path)]
		}

		callerParts := strings.Split(runtime.FuncForPC(pc).Name(), "/")
		caller := callerParts[len(callerParts)-1]
		fmt.Printf("  %16s  %s:%d\n", path, caller, line)
	}

	fmt.Println("Why things fail: your stack is to deep for sanity...")
}

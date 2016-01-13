package utils

type FailedAssertion string

func (fa FailedAssertion) Error() string {
	return string(fa)
}

func Assert(truthy bool, msg string) {
	if !truthy {
		panic(FailedAssertion(msg))
	}
}

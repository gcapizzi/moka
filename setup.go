package moka

type FailHandler func(message string, callerSkip ...int)

var globalFailHandler FailHandler

func RegisterDoublesFailHandler(failHandler FailHandler) {
	globalFailHandler = failHandler
}

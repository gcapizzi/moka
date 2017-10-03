package moka

type FailHandler func(message string, callerSkip ...int)

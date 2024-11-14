package model

type HandleState struct{}

type HandleError struct {
	// Message is a message that will be shown to the chat user.
	Message       string
	InternalError error
}

func (he *HandleError) Error() string {
	return he.InternalError.Error()
}

type HandleResult struct {
	Responses   []string
	RepeatTimes int
}

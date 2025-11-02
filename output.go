package iac

import "time"

type Output chan Message

func (output Output) PushOK(content []byte) {
	output <- Message{
		Type:      MessageOK,
		Timestamp: time.Now(),
		Content:   content,
	}
}

func (output Output) PushError(content []byte) {
	output <- Message{
		Type:      MessageError,
		Timestamp: time.Now(),
		Content:   content,
	}
}

func FailedOutput(err error) Output {
	output := make(Output, 1)
	go func() {
		defer close(output)
		output.PushError([]byte(err.Error()))
	}()
	return output
}

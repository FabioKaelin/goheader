package logger

import (
	"crypto/rand"
	"fmt"
)

type (
	logger struct {
		requestID string
	}

	loggerInterface interface {
		Println(a ...any)
		Print(a ...any)
		PrintF(format string, a ...any)
		GetRequestID() string
	}
)

func CreateLogger() loggerInterface {
	requestID := randomString(2)
	return &logger{requestID: requestID}
}

func (l *logger) Println(a ...any) {
	a = append([]any{l.requestID}, a...)
	fmt.Println(a...)
}

func (l *logger) Print(a ...any) {
	a = append([]any{l.requestID}, a...)
	fmt.Print(a...)
}

func (l *logger) PrintF(format string, a ...any) {
	format = l.requestID + " " + format
	fmt.Printf(format, a...)
}

func randomString(length int) string {
	b := make([]byte, length+2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[2 : length+2]
}

func (l *logger) GetRequestID() string {
	return l.requestID
}

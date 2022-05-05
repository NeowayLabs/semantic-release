package time

import (
	"time"
)

type Logger interface {
	Info(s string, args ...interface{})
}

type timeControl struct {
	log Logger
}

// PrintElapsedTime can be used to measure the functions elapsed time.
// Use defer PrintElapsedTime("functionName")() at the beginning of the functions
// Args:
//  	functionName (string): Name of the function to calculate the elapsed time.
func (t *timeControl) PrintElapsedTime(functionName string) func() {
	start := time.Now()
	return func() {
		t.log.Info("%s took %v\n\n", functionName, time.Since(start))
	}
}

func New(log Logger) *timeControl {
	return &timeControl{
		log: log,
	}
}

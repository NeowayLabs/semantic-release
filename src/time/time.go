package time

import (
	"log"
	"time"
)

// GetElapsedTime can be used to measure the functions elapsed time.
// Use defer GetElapsedTime("functionName")() at the beginning of the functions
// Args:
//  	what (string): Name of the function to calculate the elapsed time.
func GetElapsedTime(what string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s took %v\n\n", what, time.Since(start))
	}
}

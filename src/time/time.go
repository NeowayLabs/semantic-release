package time

import (
	"log"
	"time"
)

// PrintElapsedTime can be used to measure the functions elapsed time.
// Use defer GetElapsedTime("functionName")() at the beginning of the functions
// Args:
//  	functionName (string): Name of the function to calculate the elapsed time.
func PrintElapsedTime(functionName string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s took %v\n\n", functionName, time.Since(start))
	}
}

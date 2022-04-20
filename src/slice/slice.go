package slice

// IsStringInSlice aims to find a string inside a slice of strings.
// It requires a full match.
// Args:
// 		value (string): String value to find.
// 		slice ([]string): Slice containing strings.
// Returns:
// 		bool: True when found, otherwise false.
func IsStringInSlice(value string, slice []string) bool {
	for i := range slice {
		if slice[i] == value {
			return true
		}
	}
	return false
}

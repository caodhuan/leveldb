package utilties

import "fmt"

// Append a human-readable printout of "value" to *str.
// Escapes any non-printable characters found in "value".
func AppendEscapedStringTo(str *string, value string) {
	totalStr := make([]byte, len(*str) + len(value))
	copy(totalStr, *str)

	for _, char := range value {
		if char >= ' ' && char <= '~' {
			totalStr = append(totalStr, byte(char))
		} else {
			buf := fmt.Sprintf("\\x%02x", char & 0xff)
			copy(totalStr, buf)
		}
	}
}


// Return a human-readable version of "value".
// Escapes any non-printable characters found in "value".
func EscapeString(value string) string {
	var result string
	AppendEscapedStringTo(&result, value)
	return result
}
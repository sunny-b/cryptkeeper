package logger

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// CustomFormatter formats logs into plain text
type CustomFormatter struct{}

// Format renders a single log entry
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Start with the log message
	output := []string{entry.Message}

	// Append custom fields if any
	for key, value := range entry.Data {
		output = append(output, fmt.Sprintf("%s=%v", key, value))
	}

	// Join all parts and add a newline
	return []byte(strings.Join(output, " ") + "\n"), nil
}

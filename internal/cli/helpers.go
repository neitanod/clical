package cli

import (
	"fmt"
	"os"
)

// ANSI color codes
const (
	colorRed   = "\033[31m"
	colorReset = "\033[0m"
)

// printRedError prints an error message in red to stderr
// This is useful with SilenceUsage: true to avoid duplicate error messages
func printRedError(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "%sError: %s%s\n", colorRed, msg, colorReset)
}

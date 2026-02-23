package errutil

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"
)

// DisplayError keeps CLI display options alongside an underlying error.
type DisplayError struct {
	Err       error
	ShowMsg   bool
	ShowColor bool
	Pause     bool
}

func (e *DisplayError) Error() string {
	if e == nil || e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

// Unwrap returns the underlying error.
func (e *DisplayError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// NewDisplayError builds an error that carries ErrorExit rendering options.
func NewDisplayError(err error, showMsg bool, showColor bool, pause bool) error {
	if err == nil {
		return nil
	}
	return &DisplayError{
		Err:       err,
		ShowMsg:   showMsg,
		ShowColor: showColor,
		Pause:     pause,
	}
}

// ExitCodeError requests process termination with a specific exit code.
type ExitCodeError struct {
	Code  int
	Pause bool
}

func (e *ExitCodeError) Error() string {
	return "exit requested"
}

// NewExitCodeError builds an exit code request.
func NewExitCodeError(code int, pause bool) error {
	return &ExitCodeError{
		Code:  code,
		Pause: pause,
	}
}

// FormatError formats an error for CLI output.
func FormatError(err error) string {
	if err == nil {
		return "ERR: unknown error"
	}
	return "ERR: " + err.Error()
}

// Exit exits program
func Exit(pause bool, exitCode int) {
	if pause {
		fmt.Fprintf(os.Stdout, "Press enter to exit...")
		bufio.NewReader(os.Stdin).ReadString('\n')
	}
	os.Exit(exitCode)
}

// ErrorRedPrintln shows red string to stderr
func ErrorRedPrintln(str string) {
	color.New(color.FgRed).Fprintln(color.Error, str)
}

// StdoutGreenPrintln shows green string to stdout
func StdoutGreenPrintln(str string) {
	color.New(color.FgGreen).Fprintln(color.Output, str)
}

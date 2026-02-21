package utils

// DisplayError keeps ErrorExit display options alongside an underlying error.
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

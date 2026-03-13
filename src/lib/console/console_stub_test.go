//go:build !windows

package console

import "testing"

func TestConsoleStubFunctions(t *testing.T) {
	t.Parallel()

	if ConsoleProcNames == "" {
		t.Fatal("ConsoleProcNames is empty")
	}

	isConsole, err := IsParentConsole()
	if err != nil {
		t.Fatalf("IsParentConsole returned error: %v", err)
	}
	if !isConsole {
		t.Fatal("IsParentConsole() = false, want true")
	}

	if err := FreeConsole(); err != nil {
		t.Fatalf("FreeConsole returned error: %v", err)
	}

	AllocConsole()
	SetConsoleTitle("title")
}

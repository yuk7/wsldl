//go:build windows

package console

import (
	"errors"
	"os"
	"reflect"
	"sync"
	"syscall"
	"testing"
	"unsafe"

	ps "github.com/mitchellh/go-ps"
)

type fakeProcess struct {
	pid        int
	ppid       int
	executable string
}

func (p fakeProcess) Pid() int {
	return p.pid
}

func (p fakeProcess) PPid() int {
	return p.ppid
}

func (p fakeProcess) Executable() string {
	return p.executable
}

type fakeKernel32Proc struct {
	addr uintptr

	calls [][]uintptr
}

func (p *fakeKernel32Proc) Call(args ...uintptr) (uintptr, uintptr, error) {
	p.calls = append(p.calls, append([]uintptr(nil), args...))
	return 0, 0, nil
}

func (p *fakeKernel32Proc) Addr() uintptr {
	return p.addr
}

type fakeKernel32DLL struct {
	procs map[string]kernel32Proc
	err   error

	findProcCalls []string
}

func (d *fakeKernel32DLL) FindProc(name string) (kernel32Proc, error) {
	d.findProcCalls = append(d.findProcCalls, name)
	if d.err != nil {
		return nil, d.err
	}
	proc, ok := d.procs[name]
	if !ok {
		return nil, errors.New("proc not found")
	}
	return proc, nil
}

var consoleHooksMu sync.Mutex

func withConsoleHooksLocked(t *testing.T) {
	t.Helper()

	consoleHooksMu.Lock()

	origFindProcess := findProcessFunc
	origGetPPID := getPPIDFunc
	origLoadKernel32 := loadKernel32Func
	origGetStdHandle := getStdHandleFunc
	origUTF16PtrFromString := utf16PtrFromStringFunc
	origSyscallN := syscallNFunc

	origStdout := os.Stdout
	origStderr := os.Stderr
	origStdin := os.Stdin

	t.Cleanup(func() {
		findProcessFunc = origFindProcess
		getPPIDFunc = origGetPPID
		loadKernel32Func = origLoadKernel32
		getStdHandleFunc = origGetStdHandle
		utf16PtrFromStringFunc = origUTF16PtrFromString
		syscallNFunc = origSyscallN

		os.Stdout = origStdout
		os.Stderr = origStderr
		os.Stdin = origStdin

		consoleHooksMu.Unlock()
	})
}

func TestIsParentConsole(t *testing.T) {
	withConsoleHooksLocked(t)

	const (
		ppid      = 100
		parentPID = 200
	)

	type result struct {
		process ps.Process
		err     error
	}

	parentErr := errors.New("parent error")
	grandParentErr := errors.New("grand parent error")

	tests := []struct {
		name     string
		results  map[int]result
		want     bool
		wantErr  error
		wantPIDs []int
	}{
		{
			name: "parent matches console process case-insensitively",
			results: map[int]result{
				ppid: {process: fakeProcess{pid: ppid, ppid: parentPID, executable: "PoWeRsHeLl.ExE"}},
			},
			want:     true,
			wantPIDs: []int{ppid},
		},
		{
			name: "grand parent matches console process",
			results: map[int]result{
				ppid:      {process: fakeProcess{pid: ppid, ppid: parentPID, executable: "notepad.exe"}},
				parentPID: {process: fakeProcess{pid: parentPID, ppid: 1, executable: "cmd.exe"}},
			},
			want:     true,
			wantPIDs: []int{ppid, parentPID},
		},
		{
			name: "no parent process",
			results: map[int]result{
				ppid: {process: nil},
			},
			want:     false,
			wantPIDs: []int{ppid},
		},
		{
			name: "parent lookup returns error",
			results: map[int]result{
				ppid: {err: parentErr},
			},
			want:     false,
			wantErr:  parentErr,
			wantPIDs: []int{ppid},
		},
		{
			name: "grand parent lookup returns error",
			results: map[int]result{
				ppid:      {process: fakeProcess{pid: ppid, ppid: parentPID, executable: "notepad.exe"}},
				parentPID: {err: grandParentErr},
			},
			want:     false,
			wantErr:  grandParentErr,
			wantPIDs: []int{ppid, parentPID},
		},
		{
			name: "parent and grand parent are non console processes",
			results: map[int]result{
				ppid:      {process: fakeProcess{pid: ppid, ppid: parentPID, executable: "notepad.exe"}},
				parentPID: {process: fakeProcess{pid: parentPID, ppid: 1, executable: "explorer.exe"}},
			},
			want:     false,
			wantPIDs: []int{ppid, parentPID},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pids := make([]int, 0, 2)
			getPPIDFunc = func() int {
				return ppid
			}
			findProcessFunc = func(pid int) (ps.Process, error) {
				pids = append(pids, pid)
				r := tc.results[pid]
				return r.process, r.err
			}

			got, err := IsParentConsole()
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("err = %v, want %v", err, tc.wantErr)
			}
			if got != tc.want {
				t.Fatalf("IsParentConsole() = %v, want %v", got, tc.want)
			}
			if !reflect.DeepEqual(pids, tc.wantPIDs) {
				t.Fatalf("findProcess called with %v, want %v", pids, tc.wantPIDs)
			}
		})
	}
}

func TestFreeConsole(t *testing.T) {
	withConsoleHooksLocked(t)

	t.Run("returns load error", func(t *testing.T) {
		wantErr := errors.New("load error")
		loadKernel32Func = func() (kernel32DLL, error) {
			return nil, wantErr
		}

		err := FreeConsole()
		if !errors.Is(err, wantErr) {
			t.Fatalf("err = %v, want %v", err, wantErr)
		}
	})

	t.Run("returns find proc error", func(t *testing.T) {
		wantErr := errors.New("find proc error")
		loadKernel32Func = func() (kernel32DLL, error) {
			return &fakeKernel32DLL{err: wantErr}, nil
		}

		err := FreeConsole()
		if !errors.Is(err, wantErr) {
			t.Fatalf("err = %v, want %v", err, wantErr)
		}
	})

	t.Run("calls FreeConsole proc", func(t *testing.T) {
		proc := &fakeKernel32Proc{}
		dll := &fakeKernel32DLL{
			procs: map[string]kernel32Proc{
				"FreeConsole": proc,
			},
		}
		loadKernel32Func = func() (kernel32DLL, error) {
			return dll, nil
		}

		if err := FreeConsole(); err != nil {
			t.Fatalf("FreeConsole returned error: %v", err)
		}
		if len(dll.findProcCalls) != 1 || dll.findProcCalls[0] != "FreeConsole" {
			t.Fatalf("FindProc calls = %v, want [FreeConsole]", dll.findProcCalls)
		}
		if len(proc.calls) != 1 {
			t.Fatalf("FreeConsole proc call count = %d, want 1", len(proc.calls))
		}
	})
}

func TestAllocConsole(t *testing.T) {
	withConsoleHooksLocked(t)

	t.Run("returns when load fails", func(t *testing.T) {
		loadKernel32Func = func() (kernel32DLL, error) {
			return nil, errors.New("load error")
		}

		called := false
		getStdHandleFunc = func(n int) (syscall.Handle, error) {
			called = true
			return 0, nil
		}

		AllocConsole()
		if called {
			t.Fatal("GetStdHandle should not be called when load fails")
		}
	})

	t.Run("returns when AllocConsole proc lookup fails", func(t *testing.T) {
		loadKernel32Func = func() (kernel32DLL, error) {
			return &fakeKernel32DLL{err: errors.New("find proc error")}, nil
		}

		called := false
		getStdHandleFunc = func(n int) (syscall.Handle, error) {
			called = true
			return 0, nil
		}

		AllocConsole()
		if called {
			t.Fatal("GetStdHandle should not be called when FindProc fails")
		}
	})

	t.Run("calls AllocConsole and rebinds std handles", func(t *testing.T) {
		proc := &fakeKernel32Proc{}
		dll := &fakeKernel32DLL{
			procs: map[string]kernel32Proc{
				"AllocConsole": proc,
			},
		}
		loadKernel32Func = func() (kernel32DLL, error) {
			return dll, nil
		}

		var stdHandleCalls []int
		getStdHandleFunc = func(n int) (syscall.Handle, error) {
			stdHandleCalls = append(stdHandleCalls, n)
			switch n {
			case syscall.STD_OUTPUT_HANDLE:
				return syscall.Handle(11), nil
			case syscall.STD_ERROR_HANDLE:
				return syscall.Handle(12), nil
			case syscall.STD_INPUT_HANDLE:
				return syscall.Handle(13), nil
			default:
				t.Fatalf("unexpected handle request: %d", n)
			}
			return 0, nil
		}

		AllocConsole()

		if len(dll.findProcCalls) != 1 || dll.findProcCalls[0] != "AllocConsole" {
			t.Fatalf("FindProc calls = %v, want [AllocConsole]", dll.findProcCalls)
		}
		if len(proc.calls) != 1 {
			t.Fatalf("AllocConsole proc call count = %d, want 1", len(proc.calls))
		}

		wantStdHandleCalls := []int{
			syscall.STD_OUTPUT_HANDLE,
			syscall.STD_ERROR_HANDLE,
			syscall.STD_INPUT_HANDLE,
		}
		if !reflect.DeepEqual(stdHandleCalls, wantStdHandleCalls) {
			t.Fatalf("GetStdHandle calls = %v, want %v", stdHandleCalls, wantStdHandleCalls)
		}

		var stdoutFD uintptr
		if os.Stdout != nil {
			stdoutFD = os.Stdout.Fd()
		}
		if os.Stdout == nil || stdoutFD != uintptr(11) {
			t.Fatalf("Stdout fd = %d, want 11", stdoutFD)
		}
		var stderrFD uintptr
		if os.Stderr != nil {
			stderrFD = os.Stderr.Fd()
		}
		if os.Stderr == nil || stderrFD != uintptr(12) {
			t.Fatalf("Stderr fd = %d, want 12", stderrFD)
		}
		var stdinFD uintptr
		if os.Stdin != nil {
			stdinFD = os.Stdin.Fd()
		}
		if os.Stdin == nil || stdinFD != uintptr(13) {
			t.Fatalf("Stdin fd = %d, want 13", stdinFD)
		}
	})
}

func TestSetConsoleTitle(t *testing.T) {
	withConsoleHooksLocked(t)

	t.Run("returns when load fails", func(t *testing.T) {
		loadKernel32Func = func() (kernel32DLL, error) {
			return nil, errors.New("load error")
		}

		called := false
		syscallNFunc = func(trap uintptr, args ...uintptr) (r1, r2 uintptr, err syscall.Errno) {
			called = true
			return 0, 0, 0
		}

		SetConsoleTitle("title")
		if called {
			t.Fatal("SyscallN should not be called when load fails")
		}
	})

	t.Run("returns when proc lookup fails", func(t *testing.T) {
		loadKernel32Func = func() (kernel32DLL, error) {
			return &fakeKernel32DLL{err: errors.New("find proc error")}, nil
		}

		called := false
		syscallNFunc = func(trap uintptr, args ...uintptr) (r1, r2 uintptr, err syscall.Errno) {
			called = true
			return 0, 0, 0
		}

		SetConsoleTitle("title")
		if called {
			t.Fatal("SyscallN should not be called when proc lookup fails")
		}
	})

	t.Run("returns when UTF16 conversion fails", func(t *testing.T) {
		proc := &fakeKernel32Proc{addr: 88}
		loadKernel32Func = func() (kernel32DLL, error) {
			return &fakeKernel32DLL{
				procs: map[string]kernel32Proc{
					"SetConsoleTitleW": proc,
				},
			}, nil
		}
		utf16PtrFromStringFunc = func(s string) (*uint16, error) {
			return nil, errors.New("utf16 error")
		}

		called := false
		syscallNFunc = func(trap uintptr, args ...uintptr) (r1, r2 uintptr, err syscall.Errno) {
			called = true
			return 0, 0, 0
		}

		SetConsoleTitle("title")
		if called {
			t.Fatal("SyscallN should not be called when UTF16 conversion fails")
		}
	})

	t.Run("calls SetConsoleTitleW", func(t *testing.T) {
		proc := &fakeKernel32Proc{addr: 99}
		loadKernel32Func = func() (kernel32DLL, error) {
			return &fakeKernel32DLL{
				procs: map[string]kernel32Proc{
					"SetConsoleTitleW": proc,
				},
			}, nil
		}

		rawTitle := []uint16{'t', 'i', 't', 'l', 'e', 0}
		wantPtr := &rawTitle[0]
		utf16PtrFromStringFunc = func(s string) (*uint16, error) {
			if s != "title" {
				t.Fatalf("title = %q, want %q", s, "title")
			}
			return wantPtr, nil
		}

		var (
			gotTrap uintptr
			gotArgs []uintptr
		)
		syscallNFunc = func(trap uintptr, args ...uintptr) (r1, r2 uintptr, err syscall.Errno) {
			gotTrap = trap
			gotArgs = append([]uintptr(nil), args...)
			return 0, 0, 0
		}

		SetConsoleTitle("title")

		if gotTrap != proc.addr {
			t.Fatalf("trap = %d, want %d", gotTrap, proc.addr)
		}
		if len(gotArgs) != 1 || gotArgs[0] != uintptr(unsafe.Pointer(wantPtr)) {
			t.Fatalf("args = %v, want [%d]", gotArgs, uintptr(unsafe.Pointer(wantPtr)))
		}
	})
}

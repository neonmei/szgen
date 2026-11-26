package traces

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// readProcInfo reads what a process is running from /proc. Missing pieces are left empty.
func readProcInfo(pid int) procInfo {
	var info procInfo

	if path, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", pid)); err == nil {
		info.path = path
	}

	if comm, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid)); err == nil {
		info.comm = strings.TrimSuffix(string(comm), "\n")
	}

	if cmdline, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid)); err == nil && len(cmdline) > 0 {
		info.args = strings.Split(strings.TrimSuffix(string(cmdline), "\x00"), "\x00")
	}

	return info
}

// threadGroupID returns the pid of the process that owns the task tid.
func threadGroupID(tid int) (int, error) {
	status, err := os.ReadFile(fmt.Sprintf("/proc/%d/status", tid))
	if err != nil {
		return 0, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(status))
	for scanner.Scan() {
		if value, ok := strings.CutPrefix(scanner.Text(), "Tgid:"); ok {
			return strconv.Atoi(strings.TrimSpace(value))
		}
	}

	return 0, fmt.Errorf("no Tgid in /proc/%d/status", tid)
}

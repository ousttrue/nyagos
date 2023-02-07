//go:build !windows
// +build !windows

package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (this *CmdExe) run() (int, error) {
	args := []string{
		"/bin/sh",
		"-c",
		this.Cmdline,
	}
	cmd := exec.Cmd{
		Path:   "/bin/sh",
		Args:   args,
		Stdin:  this.Stdin,
		Stdout: this.Stdout,
		Stderr: this.Stderr,
	}
	if cmd.Stdin == nil {
		cmd.Stdin = os.Stdin
	}
	if cmd.Stdout == nil {
		cmd.Stdout = os.Stdout
	}
	if cmd.Stderr == nil {
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Start(); err != nil {
		return -1, err
	}
	if this.OnExec != nil && cmd.Process != nil {
		this.OnExec(cmd.Process.Pid)
	}
	if err := cmd.Wait(); err != nil {
		return -1, err
	}
	if this.OnDone != nil && cmd.Process != nil {
		this.OnDone(cmd.Process.Pid)
	}
	return cmd.ProcessState.ExitCode(), nil
}

func (this *Source) callBatch(tmpfile string) (int, error) {
	var cmdline strings.Builder

	cmdline.WriteByte('.')

	if fullpath, err := filepath.Abs(strings.ReplaceAll(this.Args[0], `"`, ``)); err == nil {
		fmt.Fprintf(&cmdline, ` "%s"`, fullpath)
	} else {
		cmdline.WriteByte(' ')
		cmdline.WriteString(this.Args[0])
	}

	for _, arg1 := range this.Args[1:] {
		cmdline.WriteByte(' ')
		cmdline.WriteString(arg1)
	}
	cmdline.WriteString(` ; (pwd ; env) > '`)
	cmdline.WriteString(tmpfile)
	cmdline.WriteString(`'`)

	return CmdExe{
		Cmdline: cmdline.String(),
		Stdin:   this.Stdin,
		Stdout:  this.Stdout,
		Stderr:  this.Stderr,
		Env:     this.Env,
		OnExec:  this.OnExec,
		OnDone:  this.OnDone,
	}.Run()
}

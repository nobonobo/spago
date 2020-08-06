package commands

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Command ...
type Command interface {
	Name() string
	Usage()
	Execute(args []string) error
}

var (
	register = map[string]Command{}
)

// Register ...
func Register(cmd Command) {
	register[cmd.Name()] = cmd
}

// Execute ...
func Execute(name string, args []string) error {
	cmd, ok := register[name]
	if !ok {
		return fmt.Errorf("unknown command: %q", name)
	}
	if err := cmd.Execute(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return fmt.Errorf("%s failed", name)
	}
	return nil
}

// RunCmd ...
func RunCmd(wd string, env []string, args ...string) (string, error) {
	log.Println(strings.Join(args, " "))
	buff := bytes.NewBuffer(nil)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = append(os.Environ(), env...)
	cmd.Dir = wd
	cmd.Stderr = buff
	cmd.Stdout = buff
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(buff, err.Error())
		return buff.String(), err
	}
	return buff.String(), nil
}

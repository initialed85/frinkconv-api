package frinkconv_repl

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

// spam newlines after each command so we've got something consistent to read until
const delimiter = "\n\n\n\n"

// REPL provides a thin and limited abstraction around a `frinkconv` REPL process
type REPL struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
}

// New returns a REPL ready to use
func New() (*REPL, error) {
	r := REPL{
		cmd: exec.Command("frinkconv"),
	}

	stdin, err := r.cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := r.cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	r.stdin = stdin
	r.stdout = stdout

	err = r.cmd.Start()
	if err != nil {
		return nil, err
	}

	_, err = r.stdin.Write([]byte(delimiter))
	if err != nil {
		return nil, err
	}

	_, err = r.readUntilDelimiter(delimiter)
	if err != nil {
		return nil, err
	}

	log.Printf("!!! REPL started at PID %#+v", r.cmd.Process.Pid)

	return &r, nil
}

func (r *REPL) readUntilDelimiter(delimiter string) (string, error) {
	buf := ""
	var err error

	for {
		c := make([]byte, 1)
		_, err = r.stdout.Read(c)
		if err != nil {
			return "", err
		}

		buf += string(c)

		if strings.HasSuffix(buf, delimiter) {
			break
		}
	}

	return buf, err
}

// Convert takes a sourceValue in sourceUnits and returns a destinationValue in destinationUnits
func (r *REPL) Convert(sourceValue float64, sourceUnits string, destinationUnits string) (float64, error) {
	expression := fmt.Sprintf("%v %v -> %v%v", sourceValue, sourceUnits, destinationUnits, delimiter)

	log.Printf(">>> %#+v", expression)

	_, err := r.stdin.Write([]byte(expression))
	if err != nil {
		return 0.0, err
	}

	output, err := r.readUntilDelimiter("\n\n\n\n")
	if err != nil {
		return 0.0, err
	}

	log.Printf("<<< %#+v", output)

	output = strings.TrimSpace(output)

	// TODO: there may be more errors than this; not too important though, ParseFloat will fail anyway if non-numbers are given to it
	if strings.Contains(output, "Unconvertable expression") || strings.Contains(output, "Conformance error") {
		return 0.0, fmt.Errorf(output)
	}

	// TODO: this is of more concern; a happy path that needs some parsing to get to the number- are there other similar cases
	//   presently unhandled?
	if strings.Contains(output, "(exactly ") {
		output = strings.TrimRight(strings.Join(strings.Split(output, "(exactly ")[1:], "(exactly "), ")")
	}

	destinationValue, err := strconv.ParseFloat(output, 64)
	if err != nil {
		return 0.0, err
	}

	return destinationValue, nil
}

// Close shuts down the REPL
func (r *REPL) Close() {
	if r.cmd.Process != nil {
		log.Printf("!!! killing REPL at PID %v", r.cmd.Process.Pid)
		_ = r.cmd.Process.Kill()
	}

	_ = r.stdin.Close()
	_ = r.stdout.Close()
}

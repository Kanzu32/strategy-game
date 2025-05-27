package terminal_test

import (
	"SERV/terminal"
	"testing"
)

func TestStart(t *testing.T) {
	terminal.Start()
}

func TestLog(t *testing.T) {
	terminal.Start()
	terminal.Log("test")
}

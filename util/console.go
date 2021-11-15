package util

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"os"
	"strings"
	"syscall"
)

// Prompt asks user for an input
// If password is true it means that the characters won't be echoed
func Prompt(message string, password bool) (result string, err error) {
	fmt.Print(message)
	if password {
		var readData []byte
		readData, err = term.ReadPassword(int(syscall.Stdin))
		result = string(readData)
	} else {
		result, err = bufio.NewReader(os.Stdin).ReadString('\n')
	}
	return strings.TrimSpace(result), err
}

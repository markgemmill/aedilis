package aedilis

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
)

var red func(format string, a ...interface{}) string
var yellow func(format string, a ...interface{}) string

func init() {
	red = color.New(color.FgRed).SprintfFunc()
	yellow = color.New(color.FgYellow).SprintfFunc()
}

type Console struct {
	sb strings.Builder
}

func (c *Console) Write(msg string, opt ...any) {
	output := fmt.Sprintf(msg, opt...)
	fmt.Printf("%s %s\n", yellow("[aedilis]"), output)
	c.sb.WriteString(output)
}

func (c *Console) WriteError(msg string, opt ...any) {
	output := fmt.Sprintf(msg, opt...)
	fmt.Printf("%s %s\n", red("[error]"), output)
	c.sb.WriteString(output)
}

func (c *Console) String() string {
	return c.sb.String()
}

func (c *Console) Clear() {
	c.sb.Reset()
}

func NewConsole() *Console {
	c := Console{}
	c.sb = strings.Builder{}
	return &c
}

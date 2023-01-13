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
	name   string
	prefix string
	sb     strings.Builder
}

func (c *Console) Write(msg string, opt ...any) {
	output := fmt.Sprintf(msg, opt...)
	fmt.Printf("%s %s\n", yellow(c.prefix), output)
	c.sb.WriteString(output)
	c.sb.WriteString("\n")
}

func (c *Console) WriteError(msg string, opt ...any) {
	output := fmt.Sprintf(msg, opt...)
	fmt.Printf("%s %s\n", red("[error]"), output)
	c.sb.WriteString(output)
}

func (c *Console) String() string {
	return c.sb.String()
}

func (c *Console) Clear() string {
	content := c.String()
	c.sb.Reset()
	return content
}

func NewConsole(name string) *Console {
	if name == "" {
		name = "aedilis"
	}
	c := Console{
		name:   name,
		prefix: fmt.Sprintf("[%s]", name),
	}
	c.sb = strings.Builder{}
	return &c
}

package new

import (
	"flag"
	"os"
	"strings"

	"github.com/nobonobo/spago/cmd/spago/commands"
)

// Comp ...
type Comp struct {
	*flag.FlagSet
	packageName   string
	componentName string
}

// Usage ...
func (c *Comp) Usage() {
	c.FlagSet.Usage()
}

// Execute ...
func (c *Comp) Execute(args []string) error {
	if err := c.FlagSet.Parse(args); err != nil {
		return err
	}
	c.componentName = c.FlagSet.Arg(0)
	outputName := strings.ToLower(c.componentName)
	output, err := os.Create(outputName + ".go")
	if err != nil {
		return nil
	}
	defer output.Close()
	if err := templ.Execute(output, map[string]interface{}{
		"PkgName": c.packageName,
		"Name":    c.componentName,
	}); err != nil {
		return err
	}
	html, err := os.Create(outputName + ".html")
	if err != nil {
		return nil
	}
	html.Close()
	return nil
}

func init() {
	c := &Comp{
		FlagSet: flag.NewFlagSet("new", flag.ContinueOnError),
	}
	c.FlagSet.StringVar(&c.packageName, "p", "main", "output package name")
	commands.Register(c)
}

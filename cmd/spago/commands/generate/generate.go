package generate

import (
	"bytes"
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/nobonobo/spago/cmd/spago/commands"
)

// Generate ...
type Generate struct {
	*flag.FlagSet
	outputName    string
	packageName   string
	componentName string
}

// Usage ...
func (s *Generate) Usage() {
	s.FlagSet.Usage()
}

// Execute ...
func (s *Generate) Execute(args []string) (err error) {
	if err := s.FlagSet.Parse(args); err != nil {
		return err
	}
	inputName := s.FlagSet.Arg(0)
	baseName := inputName[:len(inputName)-len(filepath.Ext(inputName))]
	name := filepath.Base(baseName)
	var input io.Reader = os.Stdin
	if len(inputName) > 0 && inputName != "-" {
		r, err := os.Open(inputName)
		if err != nil {
			return err
		}
		defer r.Close()
		input = r
	}
	if len(s.outputName) == 0 {
		if len(inputName) > 0 {
			s.outputName = baseName + "_gen.go"
		} else {
			s.outputName = "generated.go"
		}
	}
	if len(s.componentName) == 0 {
		s.componentName = strings.Title(name)
	}
	output, err := os.Create(s.outputName)
	if err != nil {
		os.Remove(s.outputName)
		return err
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	defer close(sig)
	go func() {
		_, ok := <-sig
		if ok {
			log.Println("aborted")
			os.Remove(s.outputName)
			os.Exit(1)
		}
	}()
	defer output.Close()

	log.Printf("generate: %s -> %s", inputName, s.outputName)
	converter := New()
	buffer := bytes.NewBuffer(nil)
	if err := converter.Do(buffer, input, s.packageName); err != nil {
		os.Remove(s.outputName)
		return err
	}
	if err := templ.Execute(output, map[string]interface{}{
		"PkgName":       s.packageName,
		"StdImports":    converter.StdModules,
		"Imports":       converter.ExtModules,
		"ComponentName": s.componentName,
		"Generated":     buffer.String(),
		"Properties":    converter.Properties,
		"Methods":       converter.Methods,
		"AppendCode":    converter.AppendCode,
	}); err != nil {
		os.Remove(s.outputName)
		return err
	}

	return nil
}

func init() {
	s := &Generate{
		FlagSet: flag.NewFlagSet("generate", flag.ContinueOnError),
	}
	s.FlagSet.StringVar(&s.outputName, "o", "", "output filename")
	s.FlagSet.StringVar(&s.packageName, "p", "main", "output package name")
	s.FlagSet.StringVar(&s.componentName, "c", "", "component name")
	commands.Register(s)
}

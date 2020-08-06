package deploy

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/nobonobo/spago/cmd/spago/commands"
	"github.com/nobonobo/spago/cmd/spago/commands/server"
)

var urls = map[string]string{
	"/":             "index.html",
	"/wasm_exec.js": "wasm_exec.js",
	"/main.wasm":    "main.wasm",
}

// Deploy ...
type Deploy struct {
	*flag.FlagSet
	isTinyGo bool
	destDir  string
}

// Usage ...
func (c *Deploy) Usage() {
	c.FlagSet.Usage()
}

// download ...
func (c *Deploy) download(url, dst string) error {
	log.Println(url, dst)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_ = b
	return nil
}

// Execute ...
func (c *Deploy) Execute(args []string) error {
	if err := c.FlagSet.Parse(args); err != nil {
		return err
	}
	c.destDir = c.FlagSet.Arg(0)
	d, err := os.Getwd()
	if err != nil {
		return err
	}
	tempDir, err := ioutil.TempDir("", "gobuild")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)
	dh := &server.DevHandler{
		GoCmd:        "go",
		WasmExecPath: filepath.Join(runtime.GOROOT(), "misc", "wasm", "wasm_exec.js"),
		WorkDir:      d,
		TempDir:      tempDir,
	}
	if c.isTinyGo {
		dh.GoCmd = "tinygo"
		output, err := commands.RunCmd(d, nil, "tinygo", "env", "TINYGOROOT")
		if err != nil {
			return err
		}
		dh.WasmExecPath = filepath.Join(strings.TrimSpace(output), "targets", "wasm_exec.js")
	}
	http.Handle("/", dh)
	server := httptest.NewServer(http.DefaultServeMux)
	defer server.Close()
	if err := os.MkdirAll(c.destDir, 0o755); err != nil {
		return err
	}
	for u, fname := range urls {
		if err := c.download(server.URL+u, filepath.Join(c.destDir, fname)); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	c := &Deploy{
		FlagSet: flag.NewFlagSet("deploy", flag.ContinueOnError),
	}
	c.FlagSet.BoolVar(&c.isTinyGo, "tinygo", false, "use tinygo tool chain")
	commands.Register(c)
}

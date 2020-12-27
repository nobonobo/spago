package server

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/nobonobo/spago/cmd/spago/commands"
)

// ProxyURL ...
type ProxyURL struct {
	m map[string]*url.URL
}

func (p ProxyURL) String() string {
	res := []string{}
	for k, v := range p.m {
		res = append(res, strings.Join([]string{k, v.String()}, "="))
	}
	return fmt.Sprint(res)
}

// Set ...
func (p *ProxyURL) Set(s string) error {
	params := strings.SplitN(s, "=", 2)
	if len(params) < 2 || len(params[0]) == 0 || params[0] == "/" {
		return fmt.Errorf("path must not be empty or root(/)")
	}
	u, err := url.Parse(params[1])
	if err != nil {
		return err
	}
	p.m[params[0]] = u
	return nil
}

// Server ...
type Server struct {
	*flag.FlagSet
	urls     ProxyURL
	addr     string
	isTinyGo bool
}

// Usage ...
func (s *Server) Usage() {
	s.FlagSet.Usage()
}

// Execute ...
func (s *Server) Execute(args []string) error {
	if err := s.FlagSet.Parse(args); err != nil {
		return err
	}
	log.Println("listen and serve:", s.addr)
	for p, u := range s.urls.m {
		log.Printf("proxy: %q => %q", p, u.String())
		http.Handle(p, httputil.NewSingleHostReverseProxy(u))
	}
	tempDir, err := ioutil.TempDir("", "gobuild")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	d, err := os.Getwd()
	if err != nil {
		return err
	}
	dh := &DevHandler{
		GoCmd:   "go",
		WorkDir: d,
		TempDir: tempDir,
	}
	if !s.isTinyGo {
		output, err := commands.RunCmd(d, nil, dh.GoCmd, "env", "GOROOT")
		if err != nil {
			return err
		}
		dh.WasmExecPath = filepath.Join(strings.TrimSpace(output), "misc", "wasm", "wasm_exec.js")
	} else {
		dh.GoCmd = "tinygo"
		output, err := commands.RunCmd(d, nil, dh.GoCmd, "env", "TINYGOROOT")
		if err != nil {
			return err
		}
		dh.WasmExecPath = filepath.Join(strings.TrimSpace(output), "targets", "wasm_exec.js")
	}
	http.Handle("/", dh)
	return http.Serve(l, nil)
}

func init() {
	s := &Server{
		FlagSet: flag.NewFlagSet("server", flag.ContinueOnError),
		urls:    ProxyURL{m: map[string]*url.URL{}},
	}
	s.FlagSet.StringVar(&s.addr, "addr", ":8080", "listen address")
	s.FlagSet.BoolVar(&s.isTinyGo, "tinygo", false, "use tinygo tool chain")
	s.FlagSet.Var(&s.urls, "p", `proxy url(value="path=URL")`)
	commands.Register(s)
}

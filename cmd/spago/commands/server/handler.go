package server

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/nobonobo/spago/cmd/spago/commands"
)

/*
tinygo build -target wasm -o main.wasm .
cp $(tinygo env TINYGOROOT)/targets/wasm_exec.js .

known issue:
https://github.com/tinygo-org/tinygo/blob/efdb2e852ec79486cf8556c8b53e389f1f66a0f2/targets/wasm_exec.js#L293
*/

const indexHTML = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width" />
    <script defer src="wasm_exec.js"></script><script>
      (async () => {
        const resp = await fetch('main.wasm');
        if (!resp.ok) {
          const pre = document.createElement('pre');
          pre.innerText = await resp.text();
          document.body.appendChild(pre);
          return;
        }
        const src = await resp.arrayBuffer();
        const go = new Go();
        const result = await WebAssembly.instantiate(src, go.importObject);
        go.run(result.instance);
      })();
    </script>
  </head>
  <body></body>
</html>
`

// DevHandler ...
type DevHandler struct {
	goCmd        string
	wasmExecPath string
	workDir      string
	tempDir      string
}

type counter struct {
	io.Writer
	count int64
}

func (c *counter) Write(b []byte) (int, error) {
	n, err := c.Writer.Write(b)
	c.count += int64(n)
	return n, err
}

func (c *counter) Length() int64 {
	return c.count
}

func serveFileWithGZip(w http.ResponseWriter, r *http.Request, mime string, stream io.Reader) (int64, error) {
	wr := &counter{Writer: w}
	writer := gzip.NewWriter(wr)
	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Encoding", "gzip")
	_, err := io.Copy(writer, stream)
	writer.Close()
	return wr.Length(), err
}

func serveGZipedFile(w http.ResponseWriter, r *http.Request, fpath string) error {
	f, err := os.Open(fpath + ".gz")
	if err != nil {
		return err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return err
	}
	w.Header().Set("Content-Encoding", "gzip")
	http.ServeContent(w, r, filepath.Base(fpath), info.ModTime(), f)
	return nil
}

func (dh *DevHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dir, base := path.Split(r.URL.Path)
	dir = filepath.Join(dh.workDir, filepath.Clean(dir))
	if len(base) == 0 {
		base = "index.html"
	}
	fpath := filepath.Join(dir, base)
	log.Println(r.Method, dir, base)
	if _, err := os.Stat(fpath); err == nil {
		http.ServeFile(w, r, fpath)
		return
	} else {
		log.Println(err, base)
	}
	switch base {
	case "index.html":
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, indexHTML)
		return
	case "wasm_exec.js":
		http.ServeFile(w, r, dh.wasmExecPath)
		return
	case "main.wasm":
		err := serveGZipedFile(w, r, fpath)
		if err == nil {
			return
		}
		log.Println(err)
		output := filepath.Join(dh.tempDir, "main.wasm")
		args := []string{"go", "generate", "-v", "./..."}
		out, err := commands.RunCmd(dir, []string{"GOOS=js", "GOARCH=wasm"}, args...)
		if err != nil {
			log.Println(out)
			http.Error(w, out, http.StatusInternalServerError)
			return
		}
		log.Println(out)
		args = []string{dh.goCmd, "build", "-o", output, "."}
		if dh.goCmd == "tinygo" {
			args = []string{dh.goCmd, "build", "-target", "wasm", "-o", output, "."}
		}
		out, err = commands.RunCmd(dir, []string{"GOOS=js", "GOARCH=wasm"}, args...)
		if err != nil {
			log.Println(out)
			http.Error(w, out, http.StatusInternalServerError)
			return
		}
		fp, err := os.Open(output)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		info, _ := fp.Stat()
		sz, err := serveFileWithGZip(w, r, "application/wasm", fp)
		log.Printf("build completed main.wasm(size: %.1fMiB-gzip->%.1fKiB)", float64(info.Size())/1024/1024, float64(sz)/1024)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
}

package generate

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

var prop = regexp.MustCompile(`^{{(.+)}}$`)

func replaceValue(s string) (string, bool) {
	props := prop.FindStringSubmatch(s)
	if len(props) > 1 {
		return props[1], true
	}
	return s, false
}

// Converter ...
type Converter struct {
	StdModules map[string]bool
	ExtModules map[string]string
	Properties []Property
	Methods    map[string]string
	AppendCode []string
	State      string
	importAS   string
}

// Property ...
type Property struct {
	Name string
	Type string
}

// New ...
func New() *Converter {
	return &Converter{
		StdModules: map[string]bool{},
		ExtModules: map[string]string{
			"github.com/nobonobo/spago": "",
		},
		Properties: []Property{},
		Methods:    map[string]string{},
		AppendCode: []string{},
	}
}

func (c *Converter) tag(z *html.Tokenizer) string {
	b, _ := z.TagName()
	return strings.ToLower(string(b))
}

type attr struct {
	k string
	v string
}

func parseAttrs(z *html.Tokenizer) []attr {
	res := []attr{}
	for {
		key, val, more := z.TagAttr()
		if len(key) == 0 && !more {
			break
		}
		k := string(key)
		v := string(val)
		res = append(res, attr{k: k, v: v})
		if !more {
			break
		}
	}
	//log.Print(res)
	return res
}

func (c *Converter) attrs(attrSlice []attr, indent int) string {
	tab0 := strings.Repeat("\t", indent)
	res := []string{}
	for _, attr := range attrSlice {
		k := attr.k
		v := attr.v
		if len(v) == 0 {
			v = "true"
		}
		switch {
		case strings.HasPrefix(k, "@"):
			name := k[1:]
			if len(name) > 0 {
				p, ok := replaceValue(v)
				if !ok {
					log.Panicf("event value format expected @%s=\"{{c.Func}}\": got %q", name, v)
				}
				c.Methods[name] = p
				res = append(res, fmt.Sprintf("\n%sspago.Event(%q, %s),", tab0, name, p))
			} else {
				p, ok := replaceValue(v)
				if !ok {
					log.Panicf("event value format expected @%s=\"{{raw}}\": got %q", name, v)
				}
				res = append(res, fmt.Sprintf("\n%s%s,", tab0, p))
			}
		default:
			p, ok := replaceValue(v)
			if !ok {
				res = append(res, fmt.Sprintf("\n%sspago.A(%q, %q),", tab0, k, v))
			} else {
				res = append(res, fmt.Sprintf("\n%sspago.A(%q, %s),", tab0, k, p))
			}
		}
	}
	if len(res) == 0 {
		return ""
	}
	return fmt.Sprintf("%s%s", tab0, strings.Join(res, ""))
}

func (c *Converter) generate(w io.Writer, r io.Reader) (err error) {
	indent := 1
	z := html.NewTokenizer(r)
	for err == nil {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			err = z.Err()
			if err == io.EOF {
				break
			}
			return
		case html.CommentToken:
		case html.DoctypeToken:
		case html.TextToken:
			tab := strings.Repeat("\t", indent)
			t := strings.TrimSpace(string(z.Text()))
			if len(t) > 0 {
				switch c.State {
				case "import":
					c.ExtModules[t] = c.importAS
					continue
				case "raw":
					fmt.Fprintf(w, "\n%s%s,", tab[:len(tab)-1], t)
					continue
				}
				p, ok := replaceValue(t)
				if !ok {
					fmt.Fprintf(w, "\n%sspago.T(%q),", tab, t)
				} else {
					fmt.Fprintf(w, "\n%s%s,", tab, p)
				}
			}
		case html.StartTagToken:
			tab := strings.Repeat("\t", indent)
			indent++
			attrSlice := parseAttrs(z)
			tag := c.tag(z)
			c.State = strings.ToLower(tag)
			switch c.State {
			case "import":
				as := ""
				for _, a := range attrSlice {
					if a.k == "as" {
						as = a.v
					}
				}
				c.importAS = as
				continue
			case "raw":
				continue
			}
			e := fmt.Sprintf("spago.Tag(%q, ", tag)
			a, t := c.attrs(attrSlice, indent), string(z.Text())
			if indent > 2 {
				fmt.Fprint(w, "\n")
			}
			if indent == 2 {
				tab = " "
			}
			fmt.Fprintf(w, "%s%s%s", tab, e, a)
			if len(t) > 0 {
				p, ok := replaceValue(t)
				if !ok {
					fmt.Fprintf(w, "\n%sspago.T(%q),", tab, t)
				} else {
					fmt.Fprintf(w, "\n%s%s,", tab, p)
				}
			}
		case html.SelfClosingTagToken:
			tab := strings.Repeat("\t", indent)
			indent++
			tag := c.tag(z)
			a := c.attrs(parseAttrs(z), indent)
			if indent == 2 {
				tab = " "
			}
			e := fmt.Sprintf("\n%sspago.Tag(%q, ", tab, tag)
			if len(a) > 0 {
				fmt.Fprintf(w, "\n%[1]s%s%s\n%[1]s),", tab, e, a)
			} else {
				fmt.Fprintf(w, "\n%[1]s%s),", tab, e)
			}
			indent--
		case html.EndTagToken:
			indent--
			switch c.State {
			case "import", "raw":
				c.State = ""
				continue
			}
			tab := strings.Repeat("\t", indent)
			fmt.Fprintf(w, "\n%s)", tab)
			if indent > 1 {
				fmt.Fprint(w, ",")
			}
			c.State = ""
		}
	}
	return nil
}

// Do ...
func (c *Converter) Do(output io.Writer, input io.Reader, pkg string) error {
	if err := c.generate(output, input); err != nil && err != io.EOF {
		return err
	}
	return nil
}

package mystars

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type Generater struct {
	w *bufio.Writer
}

func NewGenerater(file *os.File) *Generater {
	return &Generater{
		w: bufio.NewWriter(file),
	}
}

func (g *Generater) Title() (int, error) {
	date := time.Now().Format(time.RFC1123Z)
	return g.w.WriteString("# My Stars\nUpdated at: " + date + "\n\n")
}

func (g *Generater) Desc() (int, error) {
	c1 := "`"
	c3 := "```"

	template := `Usage:
1. Get a access token from [github](https://github.com/settings/tokens)
2. Execute:
	%sbash
	go run github.com/xiaoler/mystars/cmd/main.go -o README.md -t {your token}
	%s

You can also set the env variable %sGITHUB_TOKEN%s .

If use the Github Action, please set up %ssecrets.UPDATE_TOKEN%s.
`

	s := fmt.Sprintf(template, c3, c3, c1, c1, c1, c1)
	return g.w.WriteString(s)
}

func (g *Generater) Category(lang string) (int, error) {
	return g.w.WriteString("\n### " + lang + "\n")
}

func (g *Generater) Link(name, url string) (int, error) {
	s := fmt.Sprintf("- [%s](%s)\n", name, url)
	return g.w.WriteString(s)
}

func (g *Generater) Separator() (int, error) {
	return g.w.WriteString("\n---\n")
}

func (g *Generater) Repo(abst *Abstract) (int, error) {
	template := "- [%s](%s) - %s (`stars: %d`, `license: %s`)\n"

	s := fmt.Sprintf(template,
		abst.Name,
		abst.Url,
		abst.Desc,
		abst.StarCount,
		abst.License,
	)

	return g.w.WriteString(s)
}

func (g *Generater) Flush() {
	g.w.Flush()
}

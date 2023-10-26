package pkg

import (
	"bytes"
	"fmt"
	"go/format"
	"io/fs"
	"io/ioutil"
	"strings"
)

// HookGenerator is a generator for hooks.
type HookGenerator struct {
	buf      *bytes.Buffer
	config   *config
	savePath string
}

type config struct {
	HookName   string
	UrlPattern string
	Priority   int64
}

// NewHookGenerator returns a new HookGenerator.
func NewHookGenerator(hookName, urlPattern, savePath string, priority int64) *HookGenerator {
	return &HookGenerator{
		buf: bytes.NewBuffer(nil),
		config: &config{
			HookName:   hookName,
			UrlPattern: urlPattern,
			Priority:   priority,
		},
		savePath: savePath,
	}
}

// Generate init the hook.
func (g *HookGenerator) Generate() *HookGenerator {
	if err := hookTemplate.Execute(g.buf, g.config); err != nil {
		panic(err)
	}

	return g
}

// Format format the generated code.
func (g *HookGenerator) Format() *HookGenerator {
	formatOut, err := format.Source(g.buf.Bytes())
	if err != nil {
		panic(err)
	}
	g.buf = bytes.NewBuffer(formatOut)

	return g
}

// Flush write the generated code to file.
func (g *HookGenerator) Flush() {
	filename := fmt.Sprintf("gen_%s_hook.go", strings.ToLower(g.config.HookName))
	if err := ioutil.WriteFile(
		fmt.Sprintf("%s/%s", g.savePath, filename),
		g.buf.Bytes(),
		fs.ModePerm); err != nil {
		panic(err)
	}
	fmt.Println("[QiYin] gen file ok: ", fmt.Sprintf("%s/%s", g.savePath, filename))
}

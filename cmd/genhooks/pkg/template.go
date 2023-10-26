package pkg

import "html/template"

func checkTemplate(t string) *template.Template {
	return template.Must(template.New("").Parse(t))
}

var hookTemplate = checkTemplate(`


package hooks

import (
	"fmt"

	urltrie "github.com/IRONICBo/QiYin_BE/internal/middleware/hooks/url_trie"
	"github.com/gin-gonic/gin"
)

var _ urltrie.Hook = (*{{.HookName}})(nil)

func init() {
	urltrie.RegisterHook(&{{.HookName}}{})
	fmt.Println("RegisterHook", "Register Hook[{{.HookName}}] success...")
}

type {{.HookName}} struct {
	urltrie.Hook
}

// Patterns EDIT THIS TO YOUR OWN HOOK PATTERN
func (h {{.HookName}}) Patterns() string {
	return "{{.UrlPattern}}"
}

// Priority EDIT THIS TO YOUR OWN HOOK PRIORITY
func (h GlobalHook) Priority() int64 {
	return {{.Prority}}
}

// BeforeRun EDIT THIS TO YOUR OWN HOOK BEFORE RUN
func (h {{.HookName}}) BeforeRun(c *gin.Context) {
	c.Next()
}

// AfterRun EDIT THIS TO YOUR OWN HOOK AFTER RUN
func (h {{.HookName}}) AfterRun(c *gin.Context) {
	c.Next()
}
`)

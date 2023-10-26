package urltrie

import (
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/IRONICBo/QiYin_BE/pkg/log"
)

// Mapping for store url pattern and hook.
var hookTrie *Trie

func init() {
	hookTrie = NewTrie()
}

// RegisterHook register url & hook to trie.
func RegisterHook(hook Hook) {
	hookTrie.InsertBatch(hook.Patterns(), hook)
}

// RunHook enable hook for interceptor.
func RunHook() gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := c.Request.URL.Path

		// Get path from url
		p, err := url.Parse(raw)
		if err != nil {
			log.Errorf("Hook", "parse url error: %v", err)
		}

		path := p.Path
		hooks, ok := hookTrie.Match(path)
		if !ok {
			c.Next()

			return
		}

		// Run all before hooks
		for _, hook := range hooks {
			hook.BeforeRun(c)
		}

		// Run controllers
		c.Next()

		// Run all after hooks
		for _, hook := range hooks {
			hook.AfterRun(c)
		}
	}
}

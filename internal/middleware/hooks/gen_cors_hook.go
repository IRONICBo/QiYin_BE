package hooks

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	urltrie "github.com/IRONICBo/QiYin_BE/internal/middleware/hooks/url_trie"
)

var _ urltrie.Hook = (*CORS)(nil)

func init() {
	urltrie.RegisterHook(&CORS{})
	fmt.Println("RegisterHook", "Register Hook[CORS] success...")
}

// CORS implement urltrie.Hook.
type CORS struct {
	urltrie.Hook
}

// Patterns EDIT THIS TO YOUR OWN HOOK PATTERN.
func (h *CORS) Patterns() []string {
	return []string{
		"/*",
	}
}

// Priority EDIT THIS TO YOUR OWN HOOK PRIORITY.
func (h *CORS) Priority() int64 {
	return 0
}

// BeforeRun EDIT THIS TO YOUR OWN HOOK BEFORE RUN, DO NOT NEED USE Next() FUNCTION.
func (h *CORS) BeforeRun(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
	c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Token,X-User-Id")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, New-Token, New-Expires-At")
	c.Header("Access-Control-Allow-Credentials", "true")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}
}

// AfterRun EDIT THIS TO YOUR OWN HOOK AFTER RUN, DO NOT NEED USE Next() FUNCTION.
func (h *CORS) AfterRun(c *gin.Context) {
}

package hooks

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/IRONICBo/QiYin_BE/internal/utils"
	"github.com/gin-gonic/gin"

	urltrie "github.com/IRONICBo/QiYin_BE/internal/middleware/hooks/url_trie"
)

var _ urltrie.Hook = (*JWT)(nil)

func init() {
	urltrie.RegisterHook(&JWT{})
	fmt.Println("RegisterHook", "Register Hook[JWT] success...")
}

// JWT implement urltrie.Hook.
type JWT struct {
	urltrie.Hook
}

// Patterns EDIT THIS TO YOUR OWN HOOK PATTERN.
func (h *JWT) Patterns() []string {
	return []string{
		"/api/v1/user/*",
		"/api/v1/platform/*",
		"/api/v1/login/exit",
	}
}

// Priority EDIT THIS TO YOUR OWN HOOK PRIORITY.
func (h *JWT) Priority() int64 {
	return 0
}

// BeforeRun EDIT THIS TO YOUR OWN HOOK BEFORE RUN, DO NOT NEED USE Next() FUNCTION.
func (h *JWT) BeforeRun(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" || strings.Fields(token)[0] != "Bearer" {
		c.AbortWithStatus(http.StatusUnauthorized)

		return
	}

	claims, err := utils.ParseJwtToken(strings.Fields(token)[1])
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)

		return
	}

	// Set claims to context.
	c.Set("claims", claims)
}

// AfterRun EDIT THIS TO YOUR OWN HOOK AFTER RUN, DO NOT NEED USE Next() FUNCTION.
func (h *JWT) AfterRun(c *gin.Context) {
}

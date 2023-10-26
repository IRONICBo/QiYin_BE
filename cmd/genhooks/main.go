package main

import (
	"flag"
	"os"

	"github.com/IRONICBo/QiYin_BE/cmd/genhooks/pkg"
)

func main() {
	hookName := flag.String("name", "", "hook name")
	urlPattern := flag.String("pattern", "", "url pattern")
	priority := flag.Int64("priority", 0, "hook priority")
	savePath := flag.String("path", "../../internal/middleware/hooks", "save path")
	flag.Parse()

	if *hookName == "" || *savePath == "" {
		flag.Usage()
		os.Exit(1)
	}

	pkg.NewHookGenerator(*hookName, *urlPattern, *savePath, *priority).Generate().Format().Flush()
}

package main

import (
	"flag"

	"github.com/IRONICBo/QiYin_BE/cmd/gendao/pkg"
)

func main() {
	savePath := flag.String("path", "../../internal/dal/dao", "save path")
	flag.Parse()

	models := []interface{}{
		// Add your models here.
	}

	for _, model := range models {
		pkg.NewDaoGenerator(model, *savePath).Generate().Flush()
	}
}

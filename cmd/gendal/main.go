package main

import (
	"flag"

	"github.com/IRONICBo/QiYin_BE/internal/models/base"
	"gorm.io/gen"
)

func main() {
	outpath := flag.String("path", "../../internal/dal/gen", "output path for generated dal files")
	flag.Parse()

	g := gen.NewGenerator(gen.Config{
		OutPath:       *outpath,
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	// Add tables here to generate
	tables := []interface{}{
		// Add your tables here.
		base.UserBase{},
	}

	// Generate basic dao
	g.ApplyBasic(tables...)

	// Generate query interface with dynamic query
	// Ref: https://gorm.io/gen/
	// g.ApplyInterface()

	g.Execute()
}

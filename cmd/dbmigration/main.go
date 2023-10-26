package main

import (
	"flag"

	"github.com/IRONICBo/QiYin_BE/internal/config"
	"github.com/IRONICBo/QiYin_BE/internal/conn/db"
	"github.com/IRONICBo/QiYin_BE/internal/utils"
	"github.com/IRONICBo/QiYin_BE/pkg/log"
)

var (
	configPath string
	isDrop     bool
)

func init() {
	// arg
	flag.StringVar(&configPath, "c", "../../config.yaml", "config file path")
	flag.BoolVar(&isDrop, "d", false, "drop table if exist")
	flag.Parse()

	config.ConfigInit(configPath)
	utils.Banner()
	log.InitLogger()
	db.InitMysqlDB()
}

// migrate tables.
func main() {
	// get db instance.
	db := db.GetMysqlDB()

	// tables
	tables := []interface{}{
		// Add your tables here.
	}

	// drop tables if exist.
	if isDrop {
		for i := 0; i < len(tables); i++ {
			if db.Migrator().HasTable(&tables[i]) && db.Migrator().DropTable(&tables[i]) != nil {
				log.Panicf("QiYin Table Migration", "Drop table %T... failed", tables[i])
			}
			log.Infof("QiYin Table Migration", "Drop table %T... ok", tables[i])
		}
	}

	// migrate tables.
	for i := 0; i < len(tables); i++ {
		if db.AutoMigrate(&tables[i]) != nil {
			log.Panicf("QiYin Table Migration", "Migrate table %T... failed", tables[i])
		}
		log.Infof("QiYin Table Migration", "Migrate table %T... ok", tables[i])
	}
}

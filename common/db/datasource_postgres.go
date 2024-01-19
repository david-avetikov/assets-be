package db

import (
	"assets/common/config"
	"assets/common/util"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xlab/closer"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ds *DataSource

type DataSource struct {
	*gorm.DB
}

func GetDataSource() *DataSource {
	if ds != nil {
		return ds
	}

	uri := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		config.CoreConfig.Database.Postgres.Username,
		config.CoreConfig.Database.Postgres.Password,
		config.CoreConfig.Database.Postgres.Host,
		config.CoreConfig.Database.Postgres.Port,
		config.CoreConfig.Database.Postgres.Database,
	)

	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		log.WithError(err).Error("Couldn't connect to database postgres")
		return nil
	}

	closer.Bind(func() {
		_ = util.MustOne(db.DB()).Close()
		log.Debug("Connection to database closed")
	})

	ds = &DataSource{DB: db}

	return ds
}

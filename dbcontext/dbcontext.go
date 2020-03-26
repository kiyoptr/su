package dbcontext

import (
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	DbTypeGetter func() string
	DbConnGetter func() string
)

func Open() (db *gorm.DB, err error) {
	db, err = gorm.Open(DbTypeGetter(), DbConnGetter())
	if err != nil {
		return nil, err
	}

	db.LogMode(false)
	db.SingularTable(true)

	return
}

func OpenMem() (db *gorm.DB, err error) {
	db, err = gorm.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	db.LogMode(false)
	db.SingularTable(true)

	return
}

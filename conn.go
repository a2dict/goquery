package goquery

import (
	"sync"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

var (
	defaultDB *gorm.DB
	dbStore   sync.Map
)

type Cleanup func()

// CreateDB ...
func CreateDB(name, dialect, addr string, isDefault bool) Cleanup {
	db, err := gorm.Open(dialect, addr)
	if err != nil {
		log.Fatal("failed to init db.", err)
	}
	dbStore.Store(name, db)
	if defaultDB == nil || isDefault {
		defaultDB = db
	}
	return func() {
		db.Close()
	}
}

// GetDB ...
func GetDB(name string) *gorm.DB {
	if v, ok := dbStore.Load(name); ok {
		return v.(*gorm.DB)
	}
	return nil
}

// DB get default db...
func DB() *gorm.DB {
	return defaultDB
}
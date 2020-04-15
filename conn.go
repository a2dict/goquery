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
func CreateDB(name, dialect, addr string, isDefault bool) (*gorm.DB, Cleanup) {
	db, err := gorm.Open(dialect, addr)
	if err != nil {
		log.Fatalf("failed to create db, err:%v", err)
	}
	dbStore.Store(name, db)
	if defaultDB == nil || isDefault {
		defaultDB = db
	}
	return db, func() {
		db.Close()
	}
}

// SetDefaultDB ...
func SetDefaultDB(db *gorm.DB) {
	defaultDB = db
}

// GetDB ...
func GetDB(name string) (*gorm.DB, bool) {
	if v, ok := dbStore.Load(name); ok {
		return v.(*gorm.DB), ok
	}
	return nil, false
}

// DB get default db...
func DB() *gorm.DB {
	return defaultDB
}

package test

import (
	"github.com/a2dict/goquery"
)

// Person ...
type Person struct {
	ID      uint32 `gorm:"primary_key"`
	Name    string
	Age     int
	Profile string
	City    string
}

// Save ...
func (p *Person) Save() error {
	return goquery.DB().Save(p).Error
}

// ListPersons ...
func ListPersons() ([]Person, error) {
	var res []Person
	err := goquery.DB().Find(&res).Error
	return res, err
}

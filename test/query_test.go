package test

import (
	"fmt"
	"math/rand"
	"strings"
	"sync/atomic"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/a2dict/goquery"
	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

var idx uint32

func TestGormUtil(t *testing.T) {
	setup()
	fmt.Println("succ setup")
	ps, err := ListPersons()
	if err != nil {
		log.Fatalf("fail to find peopel. err:%v", err)
	}
	fmt.Println("person list:")
	for _, p := range ps {
		fmt.Println(p)
	}
	qfun := goquery.BuildPagedQuery(&Person{})

	Convey("Test goquery Query", t, func() {
		Convey("Test Paging", func() {
			req := &goquery.QReq{
				Page: 1,
				Size: 4,
				Q:    map[string]string{},
			}
			pw, err := qfun.Query(req)
			if err != nil {
				panic(err)
			}
			So(pw.Size, ShouldEqual, 4)
		})

		Convey("Test Order", func() {
			req := &goquery.QReq{
				Page: 1,
				Size: 4,
				Q:    map[string]string{},
				Sort: []string{"-age"},
			}
			pw, err := qfun.Query(req)
			if err != nil {
				panic(err)
			}
			ps := pw.Data.([]*Person)
			tAge := ps[0].Age
			for _, p := range ps {
				if tAge < p.Age {
					panic("Order ERR")
				}
				tAge = p.Age
			}
		})

		Convey("Test Cond GT", func() {
			req := &goquery.QReq{
				Page: 1,
				Size: 400,
				Q:    map[string]string{"age::gt": "22"},
				Sort: []string{"age"},
			}
			pw, err := qfun.Query(req)
			if err != nil {
				panic(err)
			}
			for _, p := range pw.Data.([]*Person) {
				if !(p.Age > 22) {
					panic("Cond GT ERR")
				}
			}
		})

		Convey("Test Cond GE", func() {
			req := &goquery.QReq{
				Page: 1,
				Size: 400,
				Q:    map[string]string{"age::ge": "22"},
				Sort: []string{"age"},
			}
			pw, err := qfun.Query(req)
			if err != nil {
				panic(err)
			}
			for _, p := range pw.Data.([]*Person) {
				if !(p.Age >= 22) {
					panic("Cond GE ERR")
				}
			}
		})

		Convey("Test Cond LT", func() {
			req := &goquery.QReq{
				Page: 1,
				Q:    map[string]string{"age::lt": "22"},
				Sort: []string{"-age"},
			}
			pw, err := qfun.Query(req)
			if err != nil {
				panic(err)
			}
			for _, p := range pw.Data.([]*Person) {
				if !(p.Age < 22) {
					panic("Cond LT ERR")
				}
			}
		})

		Convey("Test Cond LE", func() {
			req := &goquery.QReq{
				Page: 1,
				Q:    map[string]string{"age::le": "22"},
				Sort: []string{"-age"},
			}
			pw, err := qfun.Query(req)
			if err != nil {
				panic(err)
			}
			for _, p := range pw.Data.([]*Person) {
				if !(p.Age <= 22) {
					panic("Cond LE ERR")
				}
			}
		})

		Convey("Test Cond LE With Single `:`", func() {
			req := &goquery.QReq{
				Page: 1,
				Q:    map[string]string{"age:le": "22"},
				Sort: []string{"-age"},
			}
			pw, err := qfun.Query(req)
			if err != nil {
				panic(err)
			}
			for _, p := range pw.Data.([]*Person) {
				if !(p.Age <= 22) {
					panic("Cond LE ERR")
				}
			}
		})

		Convey("Test Cond LIKE", func() {
			req := &goquery.QReq{
				Page: 1,
				Q:    map[string]string{"profile::like": "%Batman%"},
				Sort: []string{},
			}
			pw, err := qfun.Query(req)
			if err != nil {
				panic(err)
			}
			for _, p := range pw.Data.([]*Person) {
				if !(strings.Contains(p.Profile, "Batman")) {
					panic("Cond LIKE ERR")
				}
			}
		})

		Convey("Test Cond ILIKE", func() {
			req := &goquery.QReq{
				Page: 1,
				Q:    map[string]string{"profile::like": "%batman%"},
				Sort: []string{},
			}
			pw, err := qfun.Query(req)
			if err != nil {
				panic(err)
			}
			for _, p := range pw.Data.([]*Person) {
				if !(strings.Contains(p.Profile, "Batman")) {
					panic("Cond ILIKE ERR")
				}
			}
		})

		Convey("Test Cond IN", func() {
			req := &goquery.QReq{
				Page: 1,
				Q:    map[string]string{"age::in": `["18","19","20"]`},
				Sort: []string{"-age"},
			}
			pw, err := qfun.Query(req)
			if err != nil {
				panic(err)
			}
			for _, p := range pw.Data.([]*Person) {
				if !(p.Age == 18 || p.Age == 19 || p.Age == 20) {
					panic("Cond IN ERR")
				}
			}
		})

		Convey("Test Cond NOT_IN", func() {
			req := &goquery.QReq{
				Page: 1,
				Q:    map[string]string{"age::not_in": `["18","19","20"]`},
				Sort: []string{"-age"},
			}
			pw, err := qfun.Query(req)
			if err != nil {
				panic(err)
			}
			for _, p := range pw.Data.([]*Person) {
				if p.Age == 18 || p.Age == 19 || p.Age == 20 {
					panic("Cond NOT_IN ERR")
				}
			}
		})

	})
}

func setup() {
	goquery.CreateDB("db", "sqlite3", "/tmp/gorm.db", true)
	goquery.DB().DropTableIfExists(&Person{})
	goquery.DB().AutoMigrate(&Person{})

	profiles := []string{"I'm Batman.", "I'm Iron-Man", "I'm Sherlock"}
	citys := []string{"shenzhen", "guangzhou", "beijing", "hebei"}
	for i := 0; i < 60; i++ {
		p := Person{
			ID:      atomic.AddUint32(&idx, 1),
			Name:    randStr(12),
			Age:     rand.Intn(40),
			Profile: profiles[rand.Intn(len(profiles))],
			City:    citys[rand.Intn(len(citys))],
		}
		err := p.Save()
		if err != nil {
			log.Fatalf("fail to save. err:%v", err)
		}
	}
}

func randStr(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

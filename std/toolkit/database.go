package toolkit

import (
	"sync/atomic"

	"github.com/jinzhu/gorm"
	std "github.com/PKUJohnson/solar/std"
)

type DbClient struct {
	dbs         []*gorm.DB
	slaveCursor uint64
}

func (s *DbClient) Master() *gorm.DB {
	return s.dbs[0]
}

func (s *DbClient) Slave() *gorm.DB {
	return s.dbs[s.curSlave()]
}

func (s *DbClient) TableRead(t string) *gorm.DB {
	return s.Slave().Table(t)
}

func (s *DbClient) TableWrite(t string) *gorm.DB {
	return s.Master().Table(t).Begin()
}

func (s *DbClient) curSlave() int {
	l := len(s.dbs)
	if l < 2 {
		return 0
	}
	return int(1 + (atomic.AddUint64(&s.slaveCursor, 1) % uint64(l-1)))
}

func (s *DbClient) Close() {
	for _, v := range s.dbs {
		v.Close()
	}
}

func ConfigDatabase(c ...std.ConfigMysql) *DbClient {
	var database = &DbClient{}

	database.dbs = make([]*gorm.DB, len(c))

	var ch = make(chan bool)
	for i, v := range c {
		go func(i int, v std.ConfigMysql) {
			database.dbs[i] = CreateDB(v)
			ch <- true
		}(i, v)
	}
	for i := 0; i < len(c); i++ {
		<-ch
	}

	return database
}

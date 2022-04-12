package levelDb

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDb struct {
	DbPath  string
	Handler *leveldb.DB
}

func NewDbInstance(DbPath string) *leveldb.DB {
	db, err := leveldb.OpenFile(DbPath, nil)
	if err != nil {
		fmt.Println(err)
	}
	return db
}

// Instance 初始化
func (t *LevelDb) Instance() {
	t.Handler = NewDbInstance(t.DbPath) // 创建LevelDB数据库实例
}

// Put put
func (t *LevelDb) Put(key string, value string) (error, error) {
	return t.Handler.Put([]byte(key), []byte(value), nil), nil
}

// GetOne 获取单条数据
func (t *LevelDb) GetOne(key string) ([]byte, error) {
	return t.Handler.Get([]byte(key), nil)
}

func (t *LevelDb) HasOne(key string) (bool, error) {
	defer t.Handler.Close()
	return t.Handler.Has([]byte(key), nil)
}

// GetAll 获取全部数据
func (t *LevelDb) GetAll(callFunc func(key string, value string)) {
	iter := t.Handler.NewIterator(nil, nil)
	for iter.Next() {
		callFunc(string(iter.Key()), string(iter.Value()))
	}
	iter.Release()
}

func NewLevelDbInstance(DbPath string) *LevelDb {
	levelDb := LevelDb{DbPath: DbPath}
	levelDb.Instance()
	return &levelDb
}

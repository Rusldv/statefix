package datachain

import (
	"fmt"
	"os"

	"github.com/rusldv/statefix/utillib"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	blocksPath = "/sfx/blocks"
	itemsPath  = "/sfx/items"
)

var db *leveldb.DB
var itemDB *leveldb.DB

var err error

// SetPathData устанавливает путь к базе данных
func SetPathData(path string) { // TODO одновременное подключение к нескольким БД может создавать ошибки
	path = utillib.GetPath(path)
	// Блоки
	fmt.Println("База данных блоков:", path+blocksPath)
	db, err = leveldb.OpenFile(path+blocksPath, nil)
	if err != nil {
		fmt.Println("Не получилось подключиться к LevelDB:", err)
		os.Exit(1)
	}
	// Элементы
	fmt.Println("База данных элементов:", path+itemsPath)
	itemDB, err = leveldb.OpenFile(path+itemsPath, nil)
	if err != nil {
		fmt.Println("Не получилось подключиться к LevelDB:", err)
		os.Exit(1)
	}
}

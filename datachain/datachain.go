package datachain

import (
	"errors"
	"fmt"

	"github.com/rusldv/statefix/block"
)

// BlocksPool слайс с хешами блоков
var BlocksPool []string

// SetLastBlockHash setting last block hash
func SetLastBlockHash(hash string) error {
	err := db.Put([]byte(lastHash), []byte(hash), nil)
	if err != nil {
		return err
	}
	return nil
}

// LastBlockHash getting and returned last block from database.
func LastBlockHash() (string, error) {
	data, err := db.Get([]byte(lastHash), nil)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GetState get this state of chain.
func GetState() *ChainState {
	// Данные загружаются в структуру из последнего блока LastBlock.
	last, err := LastBlockHash()
	if err != nil {
		//fmt.Println(err)
		return nil
	}
	//fmt.Println(last)
	data, err := db.Get([]byte(last), nil)
	if err != nil {
		//fmt.Println(err)
		return nil
	}
	//fmt.Println(string(data))
	block, err := block.FromBencode(data)
	if err != nil {
		//fmt.Println(err)
		return nil
	}
	//fmt.Println(block.GetChainID())
	//fmt.Println(block.GetHeight())
	//fmt.Println(block.GetHash())
	return &ChainState{
		ChainID:    block.GetChainID(),
		LastHeight: block.GetHeight(),
		LastHash:   block.GetHash(),
	}
}

// AddBlock insert this block to chain and update target hash.
func AddBlock(block IBlock) error {
	if block.Check() != true {
		return errors.New("block checked error")
	}
	blockHash := []byte(block.GetHash())
	blockBencode := block.ToBencode()
	//fmt.Println(blockHash)
	//fmt.Println(string(blockBencode))
	err = db.Put(blockHash, blockBencode, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = SetLastBlockHash(string(blockHash))
	if err != nil {
		return err
	}
	//fmt.Println("Блок записан, hash:", string(blockHash))

	return nil
}

// GetBlock returned block by hash from database.
func GetBlock(hash string) (*block.Block, error) {
	//fmt.Println("GetBlock", hash)
	data, err := db.Get([]byte(hash), nil)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(data))
	block, err := block.FromBencode(data)
	if err != nil {
		return nil, err
	}
	return block, nil
}

// Checked проверяет целостность цепочки до блока генезиса
// По блоку генезиса определяется что это за блокчейн
func Checked(genesisHash string) bool {
	//fmt.Println("Checked:", genesisHash)
	state := GetState()
	//chain := state.ChainID
	thisHash := state.LastHash
	//fmt.Println(thisHash)
	// проверка блокчейна
	BlocksPool = make([]string, state.LastHeight+1)
	for i := state.LastHeight; i > 0; i-- {
		//fmt.Println("GetBlock height:", i)
		// checked
		thisBlock, err := GetBlock(thisHash)
		if err != nil {
			fmt.Println(err)
			return false
		}
		fmt.Println(i, "-", thisHash, "-", thisBlock.Check())
		// добавляем в слайс с помощью append - он в файле blockspool
		BlocksPool[i] = thisHash
		// Устанавливаем хеш следующего ниже блока
		thisHash = thisBlock.PrevHash
	}
	BlocksPool[0] = thisHash
	fmt.Println("Блок генезиса -", thisHash)
	// Проверяем хеш генезис блока
	if genesisHash != thisHash {
		return false
	}

	//fmt.Println(BlocksPool)

	return true
}

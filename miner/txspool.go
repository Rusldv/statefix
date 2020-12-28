package miner

import (
	"fmt"

	"github.com/rusldv/statefix/block"
	"github.com/rusldv/statefix/transaction"
)

var txsMap = make(map[string]*transaction.Transaction)

// Push добавляет транзакцию в пул
func Push(tx *transaction.Transaction) {
	_, ok := txsMap[tx.Hash]
	if !ok {
		txsMap[tx.Hash] = tx
	}
}

// Show показывает пул
func Show() {
	fmt.Println(txsMap)
}

// Flush добавляет транзакции в блок через переданный указатель
func Flush(b *block.Block) {
	for _, tx := range txsMap {
		b.AddTransaction(tx)
		delete(txsMap, tx.Hash)
	}
	// Думаю delete эффективнее экономит память чем txsMap = nil тута
}

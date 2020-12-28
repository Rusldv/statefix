package datachain

import (
	"fmt"

	"github.com/rusldv/statefix/utillib"

	"github.com/rusldv/statefix/cryptolib"

	"github.com/rusldv/statefix/block"
	"github.com/rusldv/statefix/transaction"
)

// Genesis create a genesis block.
func Genesis(chID int64) *block.Block {
	// new account
	PK, err := cryptolib.HexToPK(cryptolib.GenerateHex())
	if err != nil {
		fmt.Println(err)
		return nil
	}
	//fmt.Println(PK)
	pub := cryptolib.PublicKeyMarshalHex(&PK.PublicKey)
	//fmt.Println(pub)

	// create transaction
	tx1 := transaction.NewTransaction()
	tx1.SetType(1)
	tx1.SetBlockHash("")
	tx1.SetPublicKey(pub)

	tx1.AddField("Genesis")

	tx1.SetHash(utillib.BytesToSHA256Hex(tx1.Bytes()))
	tx1sig, err := cryptolib.SignDataASN1Hex(PK, utillib.BytesToSHA256(tx1.Bytes()))
	if err != nil {
		fmt.Println(err)
	}
	tx1.SetSig(tx1sig)

	//fmt.Println(tx1)
	//fmt.Println(tx1.Check())

	// create block and settings height
	block := block.NewBlock()
	block.SetChainID(chID)
	block.SetHeight(0)
	block.SetPrevHash("0")
	block.SetMerkleRoot("0")
	block.SetPublicKey(pub)
	block.SetVersion(1) // TODO in remove in config.json

	block.AddTransaction(tx1)
	block.SetHash(utillib.BytesToSHA256Hex(block.BytesHeader()))
	blockSig, err := cryptolib.SignDataASN1Hex(PK, utillib.BytesToSHA256(block.BytesHeader()))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	block.SetSig(blockSig)

	//fmt.Println(block)
	//fmt.Println(block.Check())

	// Поле BlockHash в транзакциях внутри блока заполняем хешем полученного блока
	block.Transactions[0].SetBlockHash(block.GetHash())

	return block
}

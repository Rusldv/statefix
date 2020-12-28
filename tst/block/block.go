package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"strings"

	"github.com/onrik/gomerkle"
	"github.com/rusldv/statefix/cryptolib"
	"github.com/rusldv/statefix/utillib"

	"github.com/rusldv/statefix/transaction"

	"github.com/rusldv/kit/fileutil"

	"github.com/rusldv/statefix/block"
	"github.com/rusldv/statefix/datachain"
)

var fname = flag.String("f", "block", "File path to block.")
var chainID = flag.Int("chid", 2, "Chain ID")
var pk = flag.String("pk", "", "File to private key.")
var txs = flag.String("txs", "", "Transaction list.")

func main() {
	flag.Parse()
	fmt.Println("block", *chainID, *pk, *txs)

	// Подключение базы данных блокчейна
	cfg := "default"
	datachain.SetPathData(cfg)

	if *pk == "" {
		fmt.Println("Error: empty private key")
		return
	}
	priv, err := fileutil.ReadFileString(*pk)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(priv)
	PK, err := cryptolib.HexToPK(priv)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(PK)

	if *txs == "" {
		fmt.Println("Transactions list empty.")
		return
	}

	// new block
	block := block.NewBlock()
	block.SetChainID(int64(*chainID))

	sl := strings.Split(*txs, ",")
	// Загрузка и проверка транзакций добавляемых в блок
	for _, v := range sl {
		// Добавить расширение .json?
		fmt.Println(v)
		data, err := fileutil.ReadFileString(v)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//fmt.Println(data)
		tx, err := transaction.FromJSON([]byte(data))
		if err != nil {
			fmt.Println(err)
			continue
		}
		//fmt.Println(tx)
		//fmt.Println(tx.Check())
		if !tx.Check() {
			fmt.Println("Non checked transaction", v)
			continue
		}
		block.AddTransaction(tx)
	}

	// Узнаем высоту цепочки и хеш предыдущего блока
	chainState := datachain.GetState()
	//fmt.Println(chainState)
	// Сначала устанавливаетя высота и хеш предыдущего блокаэ
	block.SetHeight(chainState.LastHeight + 1)
	block.SetPrevHash(chainState.LastHash)
	// Дерево Меркла
	tree := gomerkle.NewTree(sha256.New())
	for i := 0; i < len(block.Transactions); i++ {
		bHex, err := hex.DecodeString(block.Transactions[i].Hash)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//fmt.Println(bHex)
		tree.AddData(bHex)
	}
	err = tree.Generate()
	if err != nil {
		fmt.Println(err)
		return
	}
	roothex := hex.EncodeToString(tree.Root())
	//fmt.Println(roothex)
	block.SetMerkleRoot(roothex)
	// Версия
	block.SetVersion(1)
	// Публичный ключ заверителя блока
	block.SetPublicKey(cryptolib.PublicKeyMarshalHex(cryptolib.XYToPublicKey(PK.PublicKey.X, PK.PublicKey.Y)))
	// Задаем хэш этого заголовка блока
	block.SetHash(utillib.BytesToSHA256Hex(block.BytesHeader()))
	// этот хэш мы добавляем во все транзакции блока
	for i := 0; i < len(block.Transactions); i++ {
		block.Transactions[i].SetBlockHash(block.Hash)
	}
	// Подписываем блок
	sig, err := cryptolib.SignDataASN1Hex(PK, utillib.BytesToSHA256(block.BytesHeader()))
	if err != nil {
		fmt.Println(err)
		return
	}
	block.SetSig(sig)
	//
	//fmt.Println(block)
	if !block.Check() {
		fmt.Println("Error: block not valid")
		return
	}
	// Write JSON
	blockJSON := block.ToJSON()
	//fmt.Println(string(blockJSON))
	err = fileutil.WriteFileString(*fname+".json", string(blockJSON))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Writing JSON")
	// Write Bencode
	blockBencode := block.ToBencode()
	//fmt.Println(string(blockBencode))
	err = fileutil.WriteFileString(*fname+".bencode", string(blockBencode))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Writing Bencode")
}

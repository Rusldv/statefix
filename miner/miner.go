package miner

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/onrik/gomerkle"
	"github.com/rusldv/kit/fileutil"
	"github.com/rusldv/statefix/cryptolib"
	"github.com/rusldv/statefix/datachain"
	"github.com/rusldv/statefix/transaction"

	"github.com/rusldv/statefix/utillib"

	"github.com/rusldv/statefix/block"
)

// Mine all flag miner mode
var Mine bool

var mineBlock *block.Block

// NewOneTx создает первую транзакцию в блоке
func NewOneTx(cfg *utillib.Config) *transaction.Transaction {
	priv, err := fileutil.ReadFileString(cfg.AccountPath)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	//fmt.Println(priv)
	PK, err := cryptolib.HexToPK(priv)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	pub := cryptolib.PublicKeyMarshalHex(&PK.PublicKey)

	tx := transaction.NewTransaction()
	tx.AddField("This one transaction")
	tx.SetType(byte(1))
	tx.SetPublicKey(pub)
	tx.SetHash(utillib.BytesToSHA256Hex(tx.Bytes()))
	sig, err := cryptolib.SignDataASN1Hex(PK, utillib.BytesToSHA256(tx.Bytes()))
	if err != nil {
		fmt.Println(err)
	}
	tx.SetSig(sig)

	if !tx.Check() {
		fmt.Println("Non checked transaction")
		return nil
	}

	return tx
}

// NewMiningBlock создает майнющийся блок
func NewMiningBlock(cfg *utillib.Config) (*block.Block, error) {
	fmt.Println("Test - ChainID:", cfg)
	// Считываем из указанного в конфигурации файла аккаунта приватный ключ
	priv, err := fileutil.ReadFileString(cfg.AccountPath)
	if err != nil {
		return nil, err
	}
	// Создаем из него публичный ключ
	PK, err := cryptolib.HexToPK(priv)
	if err != nil {
		return nil, err
	}

	b := block.NewBlock()
	// Добавляем в блок первую транзакцию
	b.AddTransaction(NewOneTx(cfg))
	// И затем транзакции из пула
	Flush(b)
	fmt.Println(b)

	b.SetChainID(int64(cfg.ChainID))
	// Узнаем высоту цепочки и хеш предыдущего блока
	chStat := datachain.GetState()
	b.SetHeight(chStat.LastHeight + 1)
	b.SetPrevHash(chStat.LastHash)
	// Дерево Меркла
	tree := gomerkle.NewTree(sha256.New())
	for i := 0; i < len(b.Transactions); i++ {
		bHex, err := hex.DecodeString(b.Transactions[i].Hash)
		if err != nil {
			fmt.Println(err)
			continue
		}
		tree.AddData(bHex)
	}
	err = tree.Generate()
	if err != nil {
		return nil, err
	}
	roothex := hex.EncodeToString(tree.Root())
	b.SetMerkleRoot(roothex)

	// Версия
	b.SetVersion(cfg.Version)
	// Публичный ключ заверителя блока
	b.SetPublicKey(cryptolib.PublicKeyMarshalHex(cryptolib.XYToPublicKey(PK.PublicKey.X, PK.PublicKey.Y)))
	// Задаем хеш заголовка этого блока
	b.SetHash(utillib.BytesToSHA256Hex(b.BytesHeader()))
	// Этот хеш мы добавляем во все транзакции блока
	for i := 0; i < len(b.Transactions); i++ {
		b.Transactions[i].SetBlockHash(b.Hash)
	}

	// Подписываем блок
	sig, err := cryptolib.SignDataASN1Hex(PK, utillib.BytesToSHA256(b.BytesHeader()))
	if err != nil {
		return nil, err
	}
	b.SetSig(sig)
	// Проверка
	if !b.Check() {
		fmt.Println("No checking block", b.Hash)
		return nil, errors.New("No checked block")
	}

	return b, nil
}

// Start запускает майнинг блоков с заданным интервалом и по окончании передает блок в callback
func Start(cfg *utillib.Config, callback func(*block.Block)) {
	fmt.Println("Генерация блоков запущена...")
	go func() {
		for {
			time.Sleep(cfg.MinerInterval * time.Second)

			testBlock, err := NewMiningBlock(cfg)
			if err != nil {
				fmt.Println(err)
				return
			}
			// вызов калбека из второго аргумента этой функции
			callback(testBlock)
		}
	}()
}

package block

import (
	"fmt"
	"testing"

	"github.com/rusldv/statefix/cryptolib"
	"github.com/rusldv/statefix/utillib"

	"github.com/rusldv/statefix/transaction"
)

func TestNewBlock(t *testing.T) {
	// create transactions
	//fmt.Println("--- transactions")
	// tx1
	tx1 := transaction.NewTransaction()
	tx1.SetType(1)
	tx1.AddField("tx1field1")
	tx1.AddField("tx1field2")
	//fmt.Println(tx1)
	// tx2
	tx2 := transaction.NewTransaction()
	tx2.SetType(1)
	tx2.AddField("tx2field1")
	tx2.AddField("tx2field2")
	//fmt.Println(tx2)

	// create block
	b := NewBlock()
	// setting header
	b.SetChainID(2)
	b.SetHeight(1)
	//fmt.Println(b.GetHeight())
	b.SetPrevHash("testprevhash")
	b.SetMerkleRoot("testmerkleroot")
	// Устанавливаем публичный ключ
	PK, _ := cryptolib.HexToPK("71e2d237d1215f76d6c42d6a835e26ba8719da89cb1936b55c62fa1471f6b0c900d4f6823013c90a")
	b.SetPublicKey(cryptolib.PublicKeyMarshalHex(cryptolib.XYToPublicKey(PK.PublicKey.X, PK.PublicKey.Y)))
	b.SetVersion(1)
	// Хеш получаем когда все задано в поля (кроме подписи - она не учитывается при хешировании)
	b.SetHash(utillib.BytesToSHA256Hex(b.BytesHeader()))
	// Затем заголовок блока подписывается
	sig, err := cryptolib.SignDataASN1Hex(PK, utillib.BytesToSHA256(b.BytesHeader()))
	if err != nil {
		t.Error(err)
	}
	b.SetSig(sig)

	// Add transactions
	b.AddTransaction(tx1)
	b.AddTransaction(tx2)

	//fmt.Println("--- block")
	//fmt.Println(b)
	//fmt.Println(b.BytesHeader())
	/*
		blockJSON := b.ToJSON()
		fmt.Println(string(blockJSON))
		b2, err := FromJSON(blockJSON)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(b2)
	*/

	blockBencode := b.ToBencode()
	fmt.Println(string(blockBencode))
	b3, err := FromBencode(blockBencode)
	if err != nil {
		println(err)
	}
	fmt.Println(b3)
	fmt.Println("--- Check")
	fmt.Println(b3.Check())
}

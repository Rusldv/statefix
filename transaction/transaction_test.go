package transaction

import (
	"fmt"
	"testing"

	"github.com/rusldv/statefix/cryptolib"
	"github.com/rusldv/statefix/utillib"
)

func TestNewTransaction(t *testing.T) {
	tx := NewTransaction()
	// Устанавливаем тип транзакции
	tx.SetType(1)
	// Устанавливаем публичный ключ
	PK, _ := cryptolib.HexToPK("71e2d237d1215f76d6c42d6a835e26ba8719da89cb1936b55c62fa1471f6b0c900d4f6823013c90a")
	tx.SetPublicKey(cryptolib.PublicKeyMarshalHex(cryptolib.XYToPublicKey(PK.PublicKey.X, PK.PublicKey.Y)))

	tx.AddField("testfield1")
	tx.AddField("testfield2")
	tx.AddField("testfield3")

	// Устанавливаем хеш
	tx.SetHash(utillib.BytesToSHA256Hex(tx.Bytes()))
	// Проверяем
	fmt.Println(tx.Check())
}

func TestNewTransaction2(t *testing.T) {
	fmt.Println("TestNewTransaction2")
	tx := NewTransaction()
	tx.SetType(1)
	//PK, _ := cryptolib.HexToPK("71e2d237d1215f76d6c42d6a835e26ba8719da89cb1936b55c62fa1471f6b0c900d4f6823013c90a")
	PK, _ := cryptolib.HexToPK(cryptolib.GenerateHex())
	// set public key hex
	tx.SetPublicKey(cryptolib.PublicKeyMarshalHex(&PK.PublicKey))
	// add fields
	tx.AddField("field1")
	tx.AddField("field2")
	tx.AddField("field3")
	// add hash
	tx.SetHash(utillib.BytesToSHA256Hex(tx.Bytes()))
	// signed
	sig, err := cryptolib.SignDataASN1Hex(PK, utillib.BytesToSHA256(tx.Bytes()))
	if err != nil {
		t.Error(err)
	}
	tx.SetSig(sig)
	//fmt.Println(tx.Sig)

	// to JSON
	txJSON := tx.ToJSON()
	fmt.Println(string(txJSON))
	fmt.Println("---")
	tx2, _ := FromJSON(txJSON)
	fmt.Println(tx2)

	// to Bencode
	txBencode := tx.ToBencode()
	fmt.Println("Bencode (string):")
	fmt.Println(string(txBencode))
	fmt.Println("Bencode unmarshal:")
	tx3, err := FromBencode(txBencode)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tx3)
	fmt.Println(tx3.Check())

	// Verify
	x, y, err := cryptolib.PublicKeyUnmarshalHex(tx2.PublicKey)
	if err != nil {
		t.Error(err)
	}
	pub := cryptolib.XYToPublicKey(x, y)
	fmt.Println(cryptolib.VerifyDataASN1Hex(pub, tx2.Sig, utillib.BytesToSHA256(tx2.Bytes())))
}

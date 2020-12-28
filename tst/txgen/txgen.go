package main

import (
	"flag"
	"fmt"

	"github.com/rusldv/kit/fileutil"

	"github.com/rusldv/statefix/cryptolib"
	"github.com/rusldv/statefix/transaction"
	"github.com/rusldv/statefix/utillib"
)

const (
	pkHex = "98b6c810c1dd5b0cb44bea9532c30f11991fa50941acd79ce7dc5b11cc091dbd2ec7f05eca341cdd"
)

var f = flag.String("f", "tx.json", "Transaction data file of JSON.")
var gen = flag.Bool("gen", false, "New transaction generate and save.")
var check = flag.Bool("check", false, "Load and checking transaction.")

func main() {
	flag.Parse()
	fmt.Println("Filename:", *f)

	//fmt.Println(cryptolib.GenerateHex())
	//fmt.Println(pkHex)
	PK, err := cryptolib.HexToPK(pkHex)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(PK)
	pub := cryptolib.PublicKeyMarshalHex(cryptolib.XYToPublicKey(PK.PublicKey.X, PK.PublicKey.Y))
	//fmt.Println(pub)

	if *gen == true {
		fmt.Println("Generate")
		tx := transaction.NewTransaction()
		tx.SetType(1)
		tx.SetPublicKey(pub)
		tx.AddField("code1")
		tx.AddField("code2")
		// Подписываем именно хэш, а не байты
		tx.SetHash(utillib.BytesToSHA256Hex(tx.BytesText()))                             // Полю присваиваем текстовый хэш
		sig, err := cryptolib.SignDataASN1Hex(PK, utillib.BytesToSHA256(tx.BytesText())) // Подписывается байтовый хэш
		if err != nil {
			fmt.Println(err)
		}
		tx.SetSig(sig)

		txJSON := tx.ToJSON()
		fmt.Println(string(txJSON))

		if err := fileutil.WriteFileString(*f, string(txJSON)); err != nil {
			fmt.Println(err)
		}
		fmt.Println(*f, "complete!")
	}

	if *check == true {
		fmt.Println("Checking")
		// TODO Сделать отдельную ReadFile для чтения только байт
		str, err := fileutil.ReadFileString(*f)
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println(str)
		t, err := transaction.FromJSON([]byte(str))
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println(t)
		x, y, err := cryptolib.PublicKeyUnmarshalHex(t.PublicKey)
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println(x, y)
		p := cryptolib.XYToPublicKey(x, y)
		fmt.Println(p)

		vf := cryptolib.VerifyDataASN1Hex(p, t.Sig, utillib.BytesToSHA256(t.BytesText()))
		fmt.Println(vf)

		//fmt.Println(t.Bytes())

		//fmt.Println("---")
		//bts := []byte("1,2,3,4,5") // Если в виде строки, то без проблем
		//bts := []byte{0, 0, 0, 1} // Что бы массив
		//bts := []byte{1} // Что бы массив
		//bText := t.BytesText()
		//fmt.Println(string(bText))
		//fmt.Println(hex.EncodeToString(bts))
		//fmt.Println(utillib.BytesToSHA256Hex(bText))
	}
}

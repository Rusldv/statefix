package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/rusldv/statefix/transaction"
	"github.com/rusldv/statefix/utillib"

	"github.com/rusldv/statefix/cryptolib"

	"github.com/rusldv/kit/fileutil"
)

var txtype = flag.Int("t", 1, "Transaction type.")
var pk = flag.String("pk", "", "Private key.")
var f = flag.String("f", "transaction", "Transaction save path.")
var fds = flag.String("fds", "", "Fieleds values.")

func main() {
	flag.Parse()
	//fmt.Println("tx", *pk, *f, *fds)
	// Load priv
	if *pk == "" {
		fmt.Println("Private key not found. Please, setting -pk file_to_key.")

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
	pub := cryptolib.PublicKeyMarshalHex(&PK.PublicKey)
	//fmt.Println("PublicKey:", pub)

	// fields
	if *fds == "" {
		fmt.Println(`Fieled values empty. Please, setting -fds "value1,value2,value3".`)

		return
	}
	//fmt.Println(*fds)
	sl := strings.Split(*fds, ",")
	//fmt.Println(sl)

	// new transaction
	tx := transaction.NewTransaction()
	tx.SetType(byte(*txtype))
	tx.SetPublicKey(pub)
	for _, v := range sl {
		tx.AddField(v)
	}
	tx.SetHash(utillib.BytesToSHA256Hex(tx.Bytes()))
	sig, err := cryptolib.SignDataASN1Hex(PK, utillib.BytesToSHA256(tx.Bytes()))
	if err != nil {
		fmt.Println(err)
	}
	tx.SetSig(sig)

	if tx.Check() {
		// saved to JSON
		txJSON := tx.ToJSON()
		//fmt.Println(string(txJSON))
		if err := fileutil.WriteFileString(*f+".json", string(txJSON)); err != nil {
			fmt.Println(err)
		}
		fmt.Println(*f, "Writing to JSON")
		// saved to Bencode
		txBencode := tx.ToBencode()
		//fmt.Println(string(txBencode))
		if err := fileutil.WriteFileString(*f+".bencode", string(txBencode)); err != nil {
			fmt.Println(err)
		}
		fmt.Println(*f, "Writing to Bencode")
	} else {
		fmt.Println("Error transaction checked.")
	}
}

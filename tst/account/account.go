package main

import (
	"flag"
	"fmt"

	"github.com/rusldv/kit/fileutil"

	"github.com/rusldv/statefix/cryptolib"
)

var gen = flag.Bool("gen", false, "Generate new account.")
var key = flag.Bool("key", false, "Show key for account.")
var f = flag.String("f", "account.key", "Filename for account data key.")

func main() {
	flag.Parse()
	if *gen {
		priv := cryptolib.GenerateHex()
		fmt.Println(priv)
		if *f != "" {
			err := fileutil.WriteFileString(*f, priv)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Seve to", *f)
		}

		return
	}

	if *key {
		if *f == "" {
			fmt.Println("No file data for key.")
			return
		}
		fmt.Println("Load key from", *f)
		priv, err := fileutil.ReadFileString(*f)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Key for", priv)
		PK, err := cryptolib.HexToPK(priv)
		if err != nil {
			fmt.Println(err)
			return
		}
		pub := cryptolib.PublicKeyMarshalHex(&PK.PublicKey)
		fmt.Println("Public key:", pub)
		addr, err := cryptolib.HexToAddress(pub)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Address: 0x%s\n", addr)

		return
	}

	flag.PrintDefaults()
}

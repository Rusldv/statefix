package main

import (
	"flag"
	"fmt"

	"github.com/rusldv/kit/fileutil"

	"github.com/rusldv/statefix/block"
	"github.com/rusldv/statefix/datachain"
)

var chst *datachain.ChainState

var genesis = flag.Bool("genesis", false, "Generate genesis block.")
var explore = flag.Bool("explore", false, "Explore chain.")
var cnt = flag.Int("cnt", 5, "Count blocks explore.")
var add = flag.String("add", "", "Addtd block file path.")
var getb = flag.String("getblock", "", "Getting block by hash.")
var gett = flag.String("gettx", "", "Getting transaction by hash.")

func main() {
	flag.Parse()

	datachain.SetPathData("default")

	chst = datachain.GetState()
	if chst == nil {
		fmt.Println("State is not defined")
	}

	if *genesis {
		if chst != nil {
			fmt.Println("Genesis block is exists")
		} else {
			createGenesis()
		}
	}

	if len(*add) > 0 {
		addBlock(*add)
	}

	if *explore {
		exploreChain()
	}

	if len(*getb) > 0 {
		getblock(*getb)
	}

	if len(*gett) > 0 {
		gettx(*gett)
	}
}

// genesis 2617add6091b2d5115a192f4454c9132fc3af1b39563d5a7b054f4d24af21594
func createGenesis() {
	genesisBlock := datachain.Genesis()
	fmt.Println(genesisBlock.Hash)
	err := datachain.AddBlock(genesisBlock)
	if err != nil {
		fmt.Println(err)
	}
	err = datachain.SetLastBlockHash(genesisBlock.Hash)
	if err != nil {
		fmt.Println(err)
	}
}

func addBlock(add string) {
	fmt.Println("Added block", add)
	blockJSON, err := fileutil.ReadFileString(add)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(blockJSON)
	block, err := block.FromJSON([]byte(blockJSON))
	if err != nil {
		fmt.Println(err)
		return
	}
	if !block.Check() {
		fmt.Println("Error: block not valid")
		return
	}

	err = datachain.AddBlock(block)
	if err != nil {
		fmt.Println(err)
	}
	err = datachain.SetLastBlockHash(block.Hash)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Add block", block.Hash)
	fmt.Println("Complete!")
}

func exploreChain() {
	var target *block.Block
	var err error

	//fmt.Println("explore")
	//fmt.Println(chst.LastHash)
	target, err = datachain.GetBlock(chst.LastHash)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(target.Height)
	//fmt.Println(target.Hash)
	targetHash := target.Hash
	// Проверяем target.Height, если он = 0, то больше не запрашиваем предыдущий
	//fmt.Println(*cnt - 1) // Первый блок всегда получаем до цикла
	for i := 0; i <= *cnt-1; i++ {
		// Загружаем блок из БД
		b, err := datachain.GetBlock(targetHash)
		if err != nil {
			fmt.Println(err)
			continue
		}
		targetHash = b.GetPrevHash()
		//fmt.Println(i, "block", targetHash, "---")
		fmt.Println(b)
		fmt.Println("---")

		if b.GetHeight() <= 0 {
			fmt.Println("END")
			break
		}
	}
}

func getblockfn(block string) *block.Block {
	fmt.Println("Data for block", block)
	b, err := datachain.GetBlock(block)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return b
}

func getblock(block string) {
	b := getblockfn(block)
	fmt.Println(string(b.ToJSON()))
}

func gettx(tx string) {
	fmt.Println("get transaction:", tx)
	//txaddr := "6b99cbff5695951ffd3045c34247948f57e1b8218b8cd66a24ac6aa0e33002269ba683756159f961012a21d92d7db9415c380cbb0a50c89e705e9cc48eee2c9e"
	txaddr := "ee47e3f99ddb520e42421d53bdf67f275f57d7804f99b26f341a65784d5448bf9ba683756159f961012a21d92d7db9415c380cbb0a50c89e705e9cc48eee2c9e"
	fmt.Println("tx addr:", txaddr)
	blockHash := txaddr[:64]
	fmt.Println("blockHash:", blockHash)
	b := getblockfn(blockHash)
	//fmt.Println(b.Transactions)
	txHash := txaddr[64:]
	fmt.Println("txHash:", txHash)
	for i, v := range b.Transactions {
		if v.Hash == txHash {
			fmt.Println(txHash, i, "OK")
			txInBlock := b.Transactions[i]
			fmt.Println(txInBlock)
		}
	}
}

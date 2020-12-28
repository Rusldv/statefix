package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/rusldv/statefix/datachain"
)

func handler(w http.ResponseWriter, r *http.Request) {
	pth := r.URL.String()
	//fmt.Println("Request:", pth)

	if pth == "/favicon.ico" {
		return
	}

	if pth == "/" {
		fmt.Println("index response")
		st := datachain.GetState()
		chid := fmt.Sprintf("%d", st.ChainID)
		lh := fmt.Sprintf("%d", st.LastHeight)
		fmt.Fprintf(w, "Node state: ChainID: "+chid+" LastHeight: "+lh+" LastHash: "+st.LastHash)
		return
	}

	if len(pth) < 2 {
		fmt.Println("length error: wait min 2")
		fmt.Fprintf(w, "length error: wait min 2")
		return
	}

	arr := strings.Split(pth, "/")
	if len(arr) <= 2 {
		fmt.Println("args error: wait min 2")
		fmt.Fprintf(w, "length error: wait min 2")
		return
	}

	fmt.Println(arr[1], arr[2])

	switch arr[1] {
	case "q":
		switch arr[2] {
		case "test":
			fmt.Println("test")
		case "get_block":
			res, err := getBlock(arr[3])
			if err != nil {
				fmt.Println(err) // TODO отправляем в ответ ошибку что блок не найден
				return
			}
			fmt.Fprintf(w, res.Result)
		case "get_block_num":
			res, err := getBlockNum(arr[3])
			if err != nil {
				fmt.Println(err) // TODO отправляем в ответ ошибку что блок не найден
				return
			}
			fmt.Fprintf(w, res.Result)
		default:
			fmt.Println("default")
			fmt.Fprintf(w, "query is empty")
		}
	case "t":
		fmt.Println("Transaction:", arr[2])
		tx, err := getTransaction(arr[2])
		if err != nil {
			fmt.Println(err) // TODO отправляем в ответ ошибку что блок не найден
			return
		}
		//fmt.Println(tx)
		fmt.Fprintf(w, tx.Result)
	case "b":
		fmt.Println("Block:", arr[2])
		res, err := getBlock(arr[2])
		if err != nil {
			fmt.Println(err) // TODO отправляем в ответ ошибку что блок не найден
			return
		}
		//fmt.Println(res.Result)
		fmt.Fprintf(w, res.Result)
	default:
		fmt.Println("Uncnown:", arr[1])
		fmt.Fprintf(w, "Default")
	}
}

// HTTPRun запускает HTTP сервер
func HTTPRun(port string) {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+port, nil)
}

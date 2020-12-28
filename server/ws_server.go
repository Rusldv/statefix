package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	path        = "/"
	bufferSize  = 4096
	checkOrigin = true
)

// WSMsg структура принимаемого сообщения
type WSMsg struct {
	Call    string `json:"call"`
	Payload string `json:"payload"`
	Rand    string `json:"rand"`
}

// WSRes структура ответного сообщения
type WSRes struct {
	Status string `json:"status"`
	Result string `json:"result"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  bufferSize,
	WriteBufferSize: bufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return checkOrigin
	},
}

// WSRun запускает WebSocket сервер
func WSRun(port string) error {
	//fmt.Println("WSRun on port", port)
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer ws.Close()
		ws.SetCloseHandler(func(code int, text string) error {
			fmt.Println("Сработала функция закрытия соединения")
			// TODO задать константу из настроек для time и др. настройки для корректного закрытия.
			m := websocket.FormatCloseMessage(code, "")
			ws.WriteControl(websocket.CloseMessage, m, time.Now().Add(5*time.Minute))
			ws.Close()

			return nil
		})

		fmt.Println("New connection:", ws.RemoteAddr())
		// Обработка соединения
		for { // Начало
			// Получаем сообщение от клиента
			msg := new(WSMsg)
			err := ws.ReadJSON(msg)
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println("Message from", ws.RemoteAddr())

			switch msg.Call {
			case "test":
				res, err := test(msg.Payload)
				if err != nil {
					fmt.Println(err)
					break
				}
				err = ws.WriteJSON(res)
				if err != nil {
					fmt.Println(err)
				}
			case "get_block":
				res, err := getBlock(msg.Payload)
				if err != nil {
					fmt.Println(err)
					break
				}
				err = ws.WriteJSON(res)
				if err != nil {
					fmt.Println(err)
				}
			case "get_tx":
				res, err := getTransaction(msg.Payload)
				if err != nil {
					fmt.Println(err)
					break
				}
				err = ws.WriteJSON(res)
				if err != nil {
					fmt.Println(err)
				}
			case "send_tx":
				//fmt.Println(msg)

				res, err := sendTx(msg.Payload)
				if err != nil {
					fmt.Println(err)
					break
				}
				err = ws.WriteJSON(res)
				if err != nil {
					fmt.Println(err)
				}

			default:
				res, err := any(msg.Call)
				if err != nil {
					fmt.Println(err)
					break
				}
				err = ws.WriteJSON(res)
				if err != nil {
					fmt.Println(err)
				}
			}
		} // Конец
	})
	http.ListenAndServe(":"+port, nil)

	return nil
}

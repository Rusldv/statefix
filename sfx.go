package main

import (
	"flag"
	"fmt"

	"github.com/rusldv/kit/fileutil"
	"github.com/rusldv/statefix/block"
	"github.com/rusldv/statefix/connector"
	"github.com/rusldv/statefix/cryptolib"
	"github.com/rusldv/statefix/miner"

	"github.com/rusldv/statefix/datachain"

	"github.com/rusldv/statefix/server"
	"github.com/rusldv/statefix/utillib"
)

// Config file path flag.
var conf = flag.String("config", "config.json", "Initial configuration file.")
var newacc = flag.Bool("newacc", false, "Generate new account.")
var acc = flag.Bool("acc", false, "Show accaunt data.")
var keyfile = flag.String("keyfile", "account.key", "New account file name.")
var dev = flag.Bool("dev", false, "Developer mode.")
var devminer = flag.Bool("devminer", false, "Dev mining mode.")
var mainminer = flag.Bool("miner", false, "Mining mode.")
var http = flag.Bool("http", false, "Http client server.")
var ws = flag.Bool("ws", false, "WebSocket client server.")
var genesis = flag.Bool("genesis", false, "Create genesis block.")

// Configuration variable.
var cfg utillib.Config

func main() {
	flag.Parse()
	// Создание нового файла приватного ключа
	if *newacc {
		fmt.Println("Внимание! файл", *keyfile, "будет перезаписан.")
		fmt.Scanln()
		priv := cryptolib.GenerateHex()
		if *keyfile != "" {
			err := fileutil.WriteFileString(*keyfile, priv)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Данные нового аккаунта сохранены в", *keyfile)
		}
		return
	}
	if *acc {
		if *keyfile == "" {
			fmt.Println("Не указан файл аккаунта.")
			return
		}
		priv, err := fileutil.ReadFileString(*keyfile)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Приватный ключ", priv)
		PK, err := cryptolib.HexToPK(priv)
		if err != nil {
			fmt.Println(err)
			return
		}
		pub := cryptolib.PublicKeyMarshalHex(&PK.PublicKey)
		fmt.Println("Публичный ключ:", pub)
		addr, err := cryptolib.HexToAddress(pub)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Адрес аккаунта: 0x%s\n", addr)
		return
	}

	// Загрузка конфигурационных данных из заданного файла
	cfg, err := utillib.ReadConfig(*conf)
	if err != nil {
		fmt.Println("Файл конфигурации", *conf, "не обнаружен")
		return
	}
	fmt.Println("Конфигурация загружена из файла:", *conf)
	fmt.Println(cfg)
	// ИД цепочки блоков
	fmt.Println("Идентификатор блокчейна:", cfg.ChainID)
	chID := cfg.ChainID
	switch chID {
	case 1:
		fmt.Println("Подключение к MainNet")
	case 2:
		fmt.Println("Подключение к TestNet")
	case 3:
		fmt.Println("Подключение к Local")
	default:
		fmt.Println("Идентификатор блокчейна не обнаружен. Подключен к MainNet")
		chID = 1
	}

	// Подключение базы данных блокчейна
	fmt.Println("Путь к дирректории данных блокчейна:", cfg.DataPath)
	datachain.SetPathData(cfg.DataPath)
	// Проверяем наличие данных блокчейна в заданной дирректории
	chainState := datachain.GetState()
	if chainState == nil {
		// Проверяем конфиг на наличие хэша блока генезиса
		if len(cfg.GenesisHash) <= 0 {
			fmt.Println("Отсутствуют данные блокчейна. Чтобы сгенерировать новый введите флаг -genesis")
			if *genesis {
				fmt.Println("Создаем блок генезиса...")
				// Genesis
				genesis := datachain.Genesis(chID)
				//fmt.Println("ChainID:", genesis.ChainID)
				err := datachain.AddBlock(genesis)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("Блок генезиса создан.")
				fmt.Println("Сохраните его хэш в конфигурации вашего блокчейна:")
				fmt.Println(genesis.Hash)
			}
		}
	} else {
		// Проверка блокчейна
		fmt.Println("Проверка блокчейна...")
		fmt.Println("Блок генезиса:", datachain.Checked(cfg.GenesisHash))
	}
	if len(cfg.GenesisHash) <= 0 {
		return
	}
	//WS OR HTTP
	if *ws && *http {
		fmt.Println("Нельзя запустить одновременно WebSocket и HTTP сервер")
	} else {
		if *ws == true {
			fmt.Println("WebSocket сервер запущен на порту", cfg.WSPort)
			go func() {
				server.WSRun(cfg.WSPort)
			}()
		} else if *http == true {
			fmt.Println("HTTP сервер запущен на порту", cfg.HTTPPort)
			go func() {
				server.HTTPRun(cfg.HTTPPort)
			}()
		}
	}

	// Mining mode
	if *mainminer {
		miner.Mine = true

		go miner.Start(cfg, func(block *block.Block) {
			fmt.Println("This callback")
			//fmt.Println(block)
			// Добавляем блок в локальную цепочку блоков
			err := datachain.AddBlock(block)
			if err != nil {
				fmt.Println(err)
			}
		})

		fmt.Println("Режим майнинга включен")
	}

	// TCP
	if *dev == true {
		fmt.Println("Режим разработчика: для выхода нажмите кнопку ВВОД")
		// тут вызываем майнер для режима разработчика (флаг -miner работает только совместно с -dev)
		if *devminer == true {
			miner.Start(cfg, func(block *block.Block) {
				fmt.Println("This callback")
				fmt.Println(block)
				// Добавляем блок в локальную цепочку блоков
				err := datachain.AddBlock(block)
				if err != nil {
					fmt.Println(err)
				}
			})
		}
		fmt.Scanln()
	}
	if *dev == false {
		// Load data from config file.
		if len(*conf) > 0 {
			// Запускаем коннектор для синхронизации блокчейна
			connector.Start(cfg) // TODO добавим калбек complete func
			// в complete вызываем если в настройках указано так miner.Start(time.Second * 30, flushed func)
			// во flushed вызываем consensus.Calculate() и если он возвращает true рассылаем блок пирам (соседним нодам)

			// Запуск TCP сервера
			fmt.Println("TCP порт:", cfg.TCPPort)
			if cfg.TCPPort != "" {
				fmt.Println("TCP сервер запускается...")

				err := server.TCPStart(cfg)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("Сервер остановлен. Программа завершена.")
			}
		}
	}
}

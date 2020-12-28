package protocol

const (
	// ErrorState состояние ошибки
	ErrorState = iota
	// GetState запрос состояния ноды
	GetState
	// GetVersion запрс версии
	GetVersion
	// GetHeight запрос высоты блокчейна
	GetHeight
	// GetLastHash запрос хеша последнего блока
	GetLastHash
	// GetBlocks запрос новых блоков
	GetBlock
	// State данные о состоянии ноды к которой выполнен запрос строкой через запятую
	State
	// OK данные успешно обработаны
	OK
	// Version инвентарь с версией
	Version
)

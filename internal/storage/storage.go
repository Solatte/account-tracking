package storage

import "database/sql"

var (
	Trade    *tradeStorage
	Listener *listenerStorage
)

func Init(client *sql.DB) {
	Trade = NewTradeStorage(client)
	Listener = NewListenerStorage(client)
}

package token

import "time"

type Token struct {
	Id         int32     `json:"id" db:"id"`
	RowToken   string    `json:"token" db:"token"`
	Expiration time.Time `json:"expiration" db:"expiration"`
	Created    time.Time `json:"created" db:"created"`
	IdAccount  int32     `json:"id_account" db:"id_account"`
}

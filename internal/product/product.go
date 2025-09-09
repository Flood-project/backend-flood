package product

import "github.com/shopspring/decimal"

type Produt struct {
	Id             int32           `json:"id" db:"id"`
	Name           string          `json:"name" db:"name"`
	Description    string          `json:"description" db:"description"`
	Id_bucha       int32           `json:"id_bucha" db:"id_bucha"`
	Id_acionamento int32           `json:"id_acionamento" db:"id_acionamento"`
	Id_base        int32           `json:"id_base" db:"id_base"`
	Capacity       int64           `json:"capacidade" db:"capacidade"`
	Value          decimal.Decimal `json:"valor" db:"valor"`
}

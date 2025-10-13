package bucha

type Bucha struct {
	ID int32 `json:"id" db:"id"`
	TipoBucha string `json:"tipobucha" db:"tipobucha" paginate:"tipobucha"`
}
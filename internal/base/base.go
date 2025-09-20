package base

type Base struct {
	ID       int32  `json:"id" db:"id"`
	TipoBase string `json:"tipobase" db:"tipobase"`
}

package acionametos

type Acionamento struct {
	ID              int32  `json:"id" db:"id"`
	TipoAcionamento string `json:"tipoacionamento" db:"tipoacionamento"`
}

package product

type Produt struct {
	Id                 int32  `json:"id" db:"id"`
	Codigo             string `json:"codigo" db:"codigo"`
	Description        string `json:"description" db:"description"`
	CapacidadeEstatica int64  `json:"capacidade_estatica" db:"capacidade_estatica"`
	CapacidadeTrabalho int64  `json:"capacidade_trabalho" db:"capacidade_trabalho"`
	Reducao            string `json:"reducao" db:"reducao"`
	AlturaBucha        int64  `json:"altura_bucha" db:"altura_bucha"`
	Curso              int64  `json:"curso" db:"curso"`
	Id_bucha           int32  `json:"id_bucha" db:"id_bucha"`
	Id_acionamento     int32  `json:"id_acionamento" db:"id_acionamento"`
	Id_base            int32  `json:"id_base" db:"id_base"`
}

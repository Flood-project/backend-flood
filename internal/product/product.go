package product

type Produt struct {
	Id                 int32  `json:"id" db:"id" paginate:"products.id"`
	Codigo             string `json:"codigo" db:"codigo" paginate:"products.codigo"`
	Description        string `json:"description" db:"description" paginate:"products.description"`
	CapacidadeEstatica int64  `json:"capacidade_estatica" db:"capacidade_estatica" paginate:"products.capacidade_estatica"`
	CapacidadeTrabalho int64  `json:"capacidade_trabalho" db:"capacidade_trabalho" paginate:"products.capacidade_trabalho"`
	Reducao            string `json:"reducao" db:"reducao" paginate:"products.reducao"`
	AlturaBucha        int64  `json:"altura_bucha" db:"altura_bucha" paginate:"products.altura_bucha"`
	Curso              int64  `json:"curso" db:"curso" paginate:"products.curso"`
	Id_bucha           int32  `json:"id_bucha" db:"id_bucha" paginate:"products.id_bucha"`
	Id_acionamento     int32  `json:"id_acionamento" db:"id_acionamento" paginate:"products.id_acionamento"`
	Id_base            int32  `json:"id_base" db:"id_base" paginate:"products.base"`
 	Ativo 		       bool   `json:"ativo" db:"ativo" paginate:"products.ativo"`
}

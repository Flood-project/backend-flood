package product

import "github.com/Flood-project/backend-flood/internal/object_store"

type ProductWithComponents struct {
	Id                 int32                   `json:"id" db:"id" paginate:"p.id"`
	Codigo             string                  `json:"codigo" db:"codigo" paginate:"p.codigo"`
	Description        string                  `json:"description" db:"description" paginate:"p.description"`
	CapacidadeEstatica int64                   `json:"capacidade_estatica" db:"capacidade_estatica" paginate:"p.capacidade_estatica"`
	CapacidadeTrabalho int64                   `json:"capacidade_trabalho" db:"capacidade_trabalho" paginate:"p.capacidade_trabalho"`
	Reducao            string                  `json:"reducao" db:"reducao" paginate:"p.reducao"`
	AlturaBucha        int64                   `json:"altura_bucha" db:"altura_bucha" paginate:"p.altura_bucha"`
	Curso              int64                   `json:"curso" db:"curso" paginate:"p.curso"`
	IdBucha            int32                   `json:"id_bucha" db:"id_bucha" paginate:"p.id_bucha"`
	TipoBucha          string                  `json:"tipo_bucha" db:"tipo_bucha" paginate:"b.tipoBucha"`
	IdAcionamento      int32                   `json:"id_acionamento" db:"id_acionamento" paginate:"p.id_acionamento"`
	TipoAcionamento    string                  `json:"tipoacionamento" db:"tipoacionamento" paginate:"a.tipoacionamento"`
	IdBase             int32                   `json:"id_base" db:"id_base" paginate:"p.base"`
	TipoBase           string                  `json:"tipobase" db:"tipobase" paginate:"bs.tipobase"`
	Images             []object_store.FileData `json:"images"`
}

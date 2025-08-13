package account_user

type Account struct {
	Id_account    int32  `json:"id_account" db:"id_account"`
	Name          string `json:"name" db:"name"`
	Email         string `json:"email" db:"email"`
	Password_hash string `json:"password_hash" db:"password_hash"`
	// Id_Group      int32  `json:"id_group" db:"id_group"`
}

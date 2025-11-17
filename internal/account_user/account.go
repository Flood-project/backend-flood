package account_user

type Account struct {
	Id_account    int32  `json:"id" db:"id"`
	Name          string `json:"name" db:"name"`
	Email         string `json:"email" db:"email"`
	Password_hash string `json:"password_hash" db:"password_hash"`
	IdUserGroup   int32  `json:"id_user_group" db:"id_user_group"`
	Active        bool   `json:"active" db:"active"`
}

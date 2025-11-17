package account_user

type AccountWithUserGroup struct {
	Id_account    int32  `json:"id" db:"id"`
	Name          string `json:"name" db:"name"`
	Email         string `json:"email" db:"email"`
	Password_hash string `json:"password_hash" db:"password_hash"`
	IdUserGroup   int32  `json:"id_user_group" db:"id_user_group"`
	GroupName     string `json:"group_name" db:"group_name"`
	Active        bool   `json:"active" db:"active"`
}
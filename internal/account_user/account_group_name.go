package account_user

type AccountGroupName struct {
	Id        int32  `json:"id" db:"id"`
	GroupName string `json:"group_name" db:"group_name"`
}

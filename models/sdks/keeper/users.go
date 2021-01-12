package modelsSdkKeeper

type Users struct {
	Data  []*User `json:"data"`
	Total int     `json:"total"`
}

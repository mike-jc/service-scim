package modelsSdkKeeper

type Error struct {
	Status      string `json:"status"`
	Error       string `json:"error"`
	Description string `json:"description"`
}

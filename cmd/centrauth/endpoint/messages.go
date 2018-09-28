package endpoint

type AuthMsg struct {
	Address   string `json:"address" binding:"required"`
	Challenge string `json:"challenge" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

type AuthTokenMsg struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

package credential

type LoginCredential struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"type"`
}

type RestoreCredential struct {
	Key      string `json:"key"`
	Username string `json:"username"`
	Password string `json:"password"`
}

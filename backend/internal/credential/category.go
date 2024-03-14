package credential

type CategoryCredential struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Private  bool   `json:"private"`
}

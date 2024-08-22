package types

type Listener struct {
	Signer       string `json:"signer"`
	IncludeError bool   `json:"include_error"`
}

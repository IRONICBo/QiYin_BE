package requestparams

type UserParams struct {
	Name      string
	Password  string
	Signature string
}

type StyleParams struct {
	Style string
}

type UserInfoParams struct {
	Name      string
	Avatar    string `json:"avatar"`
	Signature string `json:"signature,omitempty"`
	Style     int64  `json:"style"`
}

package requestparams

type QiNiuTokenParams struct {
	Ticket string `json:"ticket" binding:"required"`
}

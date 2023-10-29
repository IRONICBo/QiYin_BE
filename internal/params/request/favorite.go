package requestparams

type FavoriteParams struct {
	//UserID     string `json:userId`
	VideoId    int64 `json:"videoId"`
	ActionType int32 `json:"actionType"` // 1 点赞，-1 取消
}

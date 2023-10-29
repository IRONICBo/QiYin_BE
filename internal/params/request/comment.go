package requestparams

type CommentDelParams struct {
	CommentId int64 `json:"commentId"`
}

type CommentAddParams struct {
	VideoId     int64  `json:"videoId"`
	CommentText string `json:"commentText"`
}

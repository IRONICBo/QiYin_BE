package requestparams

type VideoUpdateParams struct {
	PlayUrl  string `json:"play_url"`
	CoverUrl string `json:"cover_url"`
	Title    string `json:"title"` // 视频名
	Desc     string `json:"desc"`
	Category int64  `json:"category"`
	Tags     string `json:"tags"`
}

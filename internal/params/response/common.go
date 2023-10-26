package responseparams

// ListPageResponse xxx list page response.
type ListPageResponse struct {
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	Total    int         `json:"total"`
	List     interface{} `json:"list"`
}

// FilePathResponse file path response.
type FilePathResponse struct {
	Path string `json:"path"`
}

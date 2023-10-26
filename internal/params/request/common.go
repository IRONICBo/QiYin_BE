package requestparams

// FILE_REQUEST_PARAMS xxx file request params.
const FILE_REQUEST_PARAMS = "file"

// ListPageParams xxx list page params.
type ListPageParams struct {
	Page     int `json:"page"     binding:"required"`
	PageSize int `json:"page_size" binding:"required"`
}

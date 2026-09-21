package dto

// PageQuery 统一分页参数。
type PageQuery struct {
	Page     int `form:"page" json:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" json:"page_size" binding:"omitempty,min=1,max=100"`
}

// Normalize 返回带默认值的页码与每页条数。
func (p *PageQuery) Normalize() (int, int) {
	page, size := p.Page, p.PageSize
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	return page, size
}

// PageResult 统一分页返回结构。
type PageResult struct {
	Items    any   `json:"items"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

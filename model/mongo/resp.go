package mongo

import "math"

type PageList struct {
	Meta struct {
		Pagination struct {
			Total       int64 `json:"total"`
			Count       int   `json:"count"`
			PerPage     int   `json:"per_page"`
			CurrentPage int   `json:"current_page"`
			TotalPages  int64 `json:"total_pages"`
		} `json:"pagination"`
	} `json:"meta"`
	Data interface{} `json:"data"`
}

func CutPageData(total int64, page, size, count int, data interface{}) (pager PageList) {
	pager.Meta.Pagination.Total = total
	pager.Meta.Pagination.TotalPages = int64(math.Ceil(float64(total) / float64(size)))
	pager.Meta.Pagination.CurrentPage = page
	pager.Meta.Pagination.PerPage = size
	pager.Meta.Pagination.Count = count
	pager.Data = data
	return
}

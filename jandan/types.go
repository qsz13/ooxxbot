package jandan

type OOXXResult struct {
	Status        string `json:"status"`
	CurrentPage   int    `json:"current_page"`
	TotalComments int    `json:"total_comments"`
	PageCount     int    `json:"page_count"`
	Count         int    `json:"count"`
	Comments      []OOXX `json:"comments"`
}

type OOXX struct {
	Content string `json:"comment_content"`
}

type PicResult struct {
	Status        string `json:"status"`
	CurrentPage   int    `json:"current_page"`
	TotalComments int    `json:"total_comments"`
	PageCount     int    `json:"page_count"`
	Count         int    `json:"count"`
	Comments      []Pic  `json:"comments"`
}

type Pic struct {
	Content string `json:"comment_content"`
}

type HotType int

const (
	OOXX_TYPE = iota
	PIC_TYPE  = iota
)

type Hot struct {
	URL     string
	Content string
	Type    HotType
}

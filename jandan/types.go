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
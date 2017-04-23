package jandan

type JandanType int

const (
	OOXX_TYPE = iota
	PIC_TYPE  = iota
)

type CommentResult struct {
	Status        string    `json:"status"`
	CurrentPage   int       `json:"current_page"`
	TotalComments int       `json:"total_comments"`
	PageCount     int       `json:"page_count"`
	Count         int       `json:"count"`
	Comments      []Comment `json:"comments"`
}

type Comment struct {
	ID      int        `json:"comment_ID,string"`
	Content string     `json:"comment_content"`
	OO      int        `json:"vote_positive,string"`
	XX      int        `json:"vote_negative,string"`
	Type    JandanType `json:"-"`
}

package social

type Data struct {
	Id      int64    `json:"id"`
	Author  string   `json:"author"`
	Content string   `json:"content"`
	Images  []string `json:"images"`
}

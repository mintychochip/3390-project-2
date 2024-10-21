package handler

type FileQueryParams struct {
	Operation string `json:"operation"`
	Column    string `json:"column"`
}

func (q *FileQueryParams) IsEmpty() bool {
	return q.Operation == "" || q.Column == ""
}

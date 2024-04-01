package response

type SearchResponse struct {
	Files  []string `json:"files"`
	Errors []Error  `json:"errors"`
}

type Error struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

package bot

type Body struct {
	Id      string `json:"group_id"`
	Message string `json:"message"`
}

type ImageData struct {
	File string `json:"file"`
}

type ImageMessage struct {
	Type string    `json:"type"`
	Data ImageData `json:"data"`
}

type ImageBody struct {
	Id      string       `json:"group_id"`
	Message ImageMessage `json:"message"`
}

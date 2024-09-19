package err

type Error struct {
	Message string `json:"error"`
}

type Errors struct {
	Errors []Error `json:"errors"`
}

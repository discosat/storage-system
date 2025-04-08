package mission

type Mission struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Bucket string `json:"bucket"`
}

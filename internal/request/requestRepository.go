package request

type RequestRepository interface {
	GetById(id string) (Request, error)
}

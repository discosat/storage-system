package mission

type MissionRepository interface {
	GetById(id int) (Mission, error)
}

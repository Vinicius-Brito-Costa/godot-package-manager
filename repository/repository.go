package repository


var repositories = map[string]Repository{
	"github": Github{},
}

type Repository interface {
	Download(name string, version string, destiny string) bool
}

func GetRepository(repo string) Repository {
	return repositories[repo]
}

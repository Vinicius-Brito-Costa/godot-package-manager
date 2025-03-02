package repository

import "godot-package-manager/gpm/file"

var repositories = map[string]Repository{
	"github": Github{},
}

type Repository interface {
	Config(plugin *file.GPPlugin) *[]byte
	Download(plugin file.GPPlugin, destiny string) bool
}

func GetRepository(repo string) Repository {
	return repositories[repo]
}

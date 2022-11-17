package options

import (
	"github.com/BurntSushi/toml"

	"gophercart/internal/repo"
)

type Options struct {
	RepoOptions *repo.Options `toml:"repo"`
}

func NewOptionsFromFile(filename string) (*Options, error) {
	opts := Options{}
	_, err := toml.DecodeFile(filename, &opts)
	if err != nil {
		return nil, err
	}
	return &opts, nil
}

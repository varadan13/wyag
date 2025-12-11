package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Repo struct {
	Worktree string
	Wyagdir  string
	Conf     string
}

func NewRepo(path string) *Repo {
	return &Repo{
		Worktree: path,
		Wyagdir:  filepath.Join(path, ".wyag"),
		Conf:     filepath.Join(path, ".wyag", "conf"),
	}
}

func InitWyag(repo *Repo) error {
	path := repo.Wyagdir
	info, err := os.Stat(path)

	if err == nil {
		if info.IsDir() {
			return fmt.Errorf("directory already exists: %s", path)
		}
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat path %s: %w", path, err)
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}

	return nil
}

type CoreConfig struct {
	RepositoryFormatVersion int  `json:"repositoryformatversion"`
	FileMode                bool `json:"filemode"`
	Bare                    bool `json:"bare"`
}

type Config struct {
	Core CoreConfig `json:"core"`
}

func InitConf(repo *Repo) error {
	path := repo.Conf
	info, err := os.Stat(path)

	if err == nil {
		if info.IsDir() {
			return fmt.Errorf("conf directory already exists: %s", path)
		} else {
			data, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read existing conf file %s: %w", path, err)
			}

			var cfg Config
			if err := json.Unmarshal(data, &cfg); err != nil {
				return fmt.Errorf("invalid conf file format %s: %w", path, err)
			}

			if cfg.Core.RepositoryFormatVersion != 0 {
				return fmt.Errorf("unsupported repository format version in %s", path)
			}

			return nil
		}
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat path %s: %w", path, err)
	}

	defaultCfg := Config{
		Core: CoreConfig{
			RepositoryFormatVersion: 0,
			FileMode:                false,
			Bare:                    false,
		},
	}

	data, err := json.MarshalIndent(defaultCfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal default config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write conf file %s: %w", path, err)
	}

	return nil
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot get working directory:", err)
		os.Exit(1)
	}

	repo := NewRepo(wd)

	err = InitWyag(repo)

	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot initialise wyag:", err)
		os.Exit(1)
	}

	err = InitConf(repo)

	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot initialise conf:", err)
		os.Exit(1)
	}

	fmt.Println("Worktree:", repo.Worktree)
	fmt.Println("Wyagdir:", repo.Wyagdir)

}

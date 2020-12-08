package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/rizinorg/rzpm/pkg/git"
	"github.com/rizinorg/rzpm/pkg/rzpackage"
)

const (
	dbSubdir = "db"
	repoName = "rz-pm-db"
)

type Database struct {
	path string
}

func New(path string) Database {
	return Database{path}
}

func (d Database) InitOrUpdate() error {
	const (
		remoteName   = "origin"
		remoteBranch = "master"
		url          = "https://github.com/rizinorg/" + repoName
	)

	log.Print("Opening " + d.path)

	repo, err := git.Open(d.path)
	if err != nil {
		// Create the repo if it does not exist
		log.Println("Creating a local database repo in " + d.path)

		repo, err = git.Init(d.path, false)
		if err != nil {
			return fmt.Errorf("could not initialize the database repo: %w", err)
		}

		log.Printf("Setting %q as master", url)

		if err := repo.AddRemote(remoteName, url); err != nil {
			return fmt.Errorf("could not add the remote: %w", err)
		}
	}

	log.Printf("Pulling the last revision from %s/%s", remoteName, remoteBranch)

	// assume origin / master
	if err := repo.Pull(remoteName, remoteBranch, nil); err != nil {
		return fmt.Errorf("could not pull the latest revision: %w", err)
	}

	return nil
}

func (d Database) Delete() error {
	return os.RemoveAll(d.path)
}

func (d Database) GetInfoFile(packageName string) (rzpackage.InfoFile, error) {
	path := filepath.Join(d.path, dbSubdir, packageName)

	return rzpackage.FromFile(path)
}

// ListAvailablePackages returns a slice of strings containing the names of all the available packages.
func (d Database) ListAvailablePackages() ([]rzpackage.Info, error) {
	dir := filepath.Join(d.path, dbSubdir)

	ifiles, err := rzpackage.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %w", dir, err)
	}

	packages := make([]rzpackage.Info, 0, len(ifiles))

	for _, p := range ifiles {
		packages = append(packages, p.Info)
	}

	return packages, nil
}

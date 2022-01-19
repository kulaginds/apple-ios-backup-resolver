package main

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

type application struct {
	SourceDir      string
	DestinationDir string

	ManifestDB *sql.DB

	DirsForCreate  map[string]struct{}
	FilesForCreate map[string]string
}

func NewApp(sourceDir, destinationDir string) *application {
	return &application{
		SourceDir:      sourceDir,
		DestinationDir: destinationDir,

		DirsForCreate:  map[string]struct{}{},
		FilesForCreate: map[string]string{},
	}
}

func (a *application) Init() error {
	var err error

	a.ManifestDB, err = sql.Open("sqlite3", fmt.Sprintf("file:%s?_journal=off", path.Join(a.SourceDir, "Manifest.db")))
	if err != nil {
		return err
	}

	return a.ManifestDB.Ping()
}

func (a *application) Run() error {
	if err := a.fillFilesStructure(); err != nil {
		return err
	}

	if err := a.createDirs(); err != nil {
		return nil
	}

	if err := a.createFiles(); err != nil {
		return nil
	}

	return nil
}

func (a *application) fillFilesStructure() error {
	rows, err := a.getFiles()

	var (
		file File
	)

	for rows.Next() {
		err = rows.Scan(
			&file.FileID,
			&file.Domain,
			&file.RelativePath,
			&file.Flags,
			&file.File,
		)
		if err != nil {
			return err
		}

		if file.RelativePath == "" {
			continue
		}

		a.saveFileStructure(file)
	}

	return nil
}

func (a *application) getFiles() (*sql.Rows, error) {
	const query = `SELECT
		fileID,
		domain,
		relativePath,
		flags,
		file
		FROM Files`

	return a.ManifestDB.Query(query)
}

func (a *application) saveFileStructure(file File) {
	dirs := make([]string, 0, 8)
	dirs = append(dirs, prepareDomainDirs(file.Domain)...)

	pathNodes := strings.Split(file.RelativePath, "/")
	for _, pathNode := range pathNodes[:len(pathNodes)-1] {
		dirs = append(dirs, pathNode)
	}

	filename := pathNodes[len(pathNodes)-1]

	if file.Flags != 1 {
		return
	}

	// if file

	dir := path.Join(dirs...)
	a.DirsForCreate[dir] = struct{}{}
	a.FilesForCreate[path.Join(dir, filename)] = file.FileID
}

var specificGroupPrefixes = []string{
	"AppDomain-",
	"AppDomainGroup-",
	"AppDomainPlugin-",
	"SysContainerDomain-",
	"SysSharedContainerDomain-",
}

func specificPrefix(path string) string {
	for _, prefix := range specificGroupPrefixes {
		if strings.HasPrefix(path, prefix) {
			return prefix
		}
	}

	return ""
}

func prepareDomainDirs(domain string) []string {
	prefix := specificPrefix(domain)
	if prefix == "" {
		return []string{domain}
	}

	return []string{prefix[:len(prefix)-1], domain[len(prefix):]}
}

func (a *application) createDirs() error {
	var err error

	for dir := range a.DirsForCreate {
		err = os.MkdirAll(path.Join(a.DestinationDir, dir), 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *application) createFiles() error {
	var (
		err            error
		sourceFilePath string
	)

	for destinationFilePath, fileID := range a.FilesForCreate {
		sourceFilePath = path.Join(a.SourceDir, fileID[:2], fileID)

		_, err = copyFile(sourceFilePath, path.Join(a.DestinationDir, destinationFilePath))
		if err != nil {
			return err
		}
	}

	return nil
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	return io.Copy(destination, source)
}

package main

type File struct {
	FileID       string
	Domain       string
	RelativePath string
	Flags        int
	File         []byte
}

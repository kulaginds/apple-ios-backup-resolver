package main

import (
	"flag"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	src := flag.String("src", "src", "iPhone backup source directory")
	dst := flag.String("dst", "dst", "iPhone backup directory with resolved files")
	flag.Parse()

	var err error

	app := NewApp(*src, *dst)

	if err = app.Init(); err != nil {
		log.Fatalln(err)
	}

	if err = app.Run(); err != nil {
		log.Fatalln(err)
	}

	log.Println("Done!")
}

package main

import (
	"context"
	"fmt"
	"imchinese/finder"
	"imchinese/repository/view"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	data, err := os.ReadFile("text.md")
	if err != nil {
		log.Fatal(err)
	}
	prettyPrint(finder.Find(string(data)))

	if err := playRepo(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func playRepo(ctx context.Context) error {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s?_foreign_keys=on", "db.sqlite")))
	if err != nil {
		return err
	}
	mr, err := view.NewRepository(db)
	if err != nil {
		return err
	}

	all, err := mr.FindAll(ctx)
	if err != nil {
		return err
	}
	for _, one := range all {
		fmt.Printf("%+v\n", one)
	}
	return nil
}

func prettyPrint(candidates []finder.Candidate) {
	for i, candidate := range candidates {
		fmt.Printf("%4d %s 「%s」 %s\n", i, candidate.Word, candidate.Line, candidate.Path)
	}
}

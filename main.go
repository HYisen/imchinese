package main

import (
	"context"
	"errors"
	"fmt"
	"imchinese/finder"
	"imchinese/repository/models"
	"imchinese/repository/view"
	"log"
	"log/slog"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	data, err := os.ReadFile("text.md")
	if err != nil {
		log.Fatal(err)
	}
	found := finder.Find(string(data))
	prettyPrint(found)

	if err := playRepo(context.Background(), found); err != nil {
		log.Fatal(err)
	}
}

func playRepo(ctx context.Context, candidates []finder.Candidate) error {
	dial := sqlite.Open(fmt.Sprintf("%s?_foreign_keys=on", "db.sqlite"))
	db, err := gorm.Open(dial, &gorm.Config{TranslateError: true})
	if err != nil {
		return err
	}
	//db = db.Debug()
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

	seen := make(map[string]bool)
	var occasion string
	for _, candidate := range candidates {
		if occasion != candidate.Occasion() {
			occasion = candidate.Occasion()
			seen = make(map[string]bool)
		}
		if seen[candidate.Word] {
			// For not plain text such as
			// > [word](#1) and [word](#2)
			// the filter here won't work.
			// So it's **pre skip**.
			slog.Warn("pre skip too fast duplicated", "word", candidate.Word, "occasion", occasion)
			continue
		}
		seen[candidate.Word] = true
		if err := mr.Save(ctx, models.Existence{
			View: &models.View{
				Name: candidate.Word,
			},
			Quote:  candidate.Quote,
			Source: candidate.Path,
			Reason: "",
			Tag:    -1,
		}); err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				slog.Warn("skip too fast duplicated", "word", candidate.Word, "occasion", occasion)
				continue
			}
			return err
		}
	}
	return nil
}

func prettyPrint(candidates []finder.Candidate) {
	for i, candidate := range candidates {
		fmt.Printf("%4d %s 「%s」 %s\n", i, candidate.Word, candidate.Quote, candidate.Path)
	}
}

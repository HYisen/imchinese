package main

import (
	"context"
	"errors"
	"flag"
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

var mode = flag.String("mode", "dump", "dump|scan")
var debug = flag.Bool("debug", false, "enable SQL debug print mode")

func main() {
	flag.Parse()

	handler, err := NewHandler(*debug)
	if err != nil {
		log.Fatal(err)
	}

	switch *mode {
	case "dump":
		handler.Dump()
	case "scan":
		handler.Scan()
	default:
		log.Fatalf("unsupported mode %s", *mode)
	}
}

type Handler struct {
	vr *view.Repository
}

func NewHandler(debug bool) (*Handler, error) {
	dial := sqlite.Open(fmt.Sprintf("%s?_foreign_keys=on", "db.sqlite"))
	db, err := gorm.Open(dial, &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, err
	}
	if debug {
		db = db.Debug()
	}

	vr, err := view.NewRepository(db)
	if err != nil {
		return nil, err
	}
	return &Handler{vr: vr}, nil
}

func (h *Handler) Dump() {
	all, err := h.vr.FindAll(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ID\tModelID\tName\tCount")
	for _, one := range all {
		fmt.Printf("%4d\t%4d\t%2d\t%s\n", one.ViewID, one.ModelID, one.Count, one.Name)
	}
}

func (h *Handler) Scan() {
	data, err := os.ReadFile("text.md")
	if err != nil {
		log.Fatal(err)
	}
	found := finder.Find(string(data))
	prettyPrint(found)

	if err := h.save(context.Background(), found); err != nil {
		log.Fatal(err)
	}
}

func (h *Handler) save(ctx context.Context, candidates []finder.Candidate) error {
	all, err := h.vr.FindAll(ctx)
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
		if err := h.vr.Save(ctx, models.Existence{
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

package models

//go:generate go tool gorm gen -i model.go -o ../generated

type Model struct {
	ID          int
	Explanation string
}

type View struct {
	ID      int
	Name    string
	ModelID int
	Model   *Model
}

type Existence struct {
	ID     int
	ViewID int
	View   *View
	Source string
	Quote  string
	Reason string
	Tag    int
	WhyNot string
}

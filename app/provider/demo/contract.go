package demo

const Key = "demo"

type Service interface {
	GetAllStudent() []Student
}

type Student struct {
	ID   int
	Name string
}

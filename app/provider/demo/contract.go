package demo

const DemoKey = "demo"

type Service interface {
	GetAllStudent() []Student
}

type Student struct {
	ID   int
	Name string
}

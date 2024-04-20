package commands

type Command interface {
	Init() error
	Run() error
}

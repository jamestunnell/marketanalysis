package models

type NewBlockFunc func() Block

type BlockRegistry interface {
	Types() []string
	Add(typ string, newBlock NewBlockFunc)
	Get(typ string) (NewBlockFunc, bool)
}

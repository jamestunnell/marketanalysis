package models

type Graph interface {
	GetBlocks() Blocks
	GetConnections() Connections
}

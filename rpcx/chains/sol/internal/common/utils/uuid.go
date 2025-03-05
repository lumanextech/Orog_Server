package utils

import (
	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node
var err error

func init() {
	// Create a new Node with a Node number of 1
	node, err = snowflake.NewNode(1)
	if err != nil {
		panic("snowflake.NewNode: " + err.Error())
	}
}

func GenerateUUID() snowflake.ID {
	return node.Generate()
}

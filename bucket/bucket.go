package bucket

import (
	"github.com/bwmarrin/snowflake"
	"math/rand"
	"time"
)

var snowflakeNode = generateSnowflakeNode()

type LengthHint interface {
	Len() int
}

func generateSnowflakeNode() *snowflake.Node {
	source := rand.NewSource(time.Now().UnixNano())
	id := source.Int63()
	high := source.Int63()
	for high > 1023 {
		high >>= 1
	}
	for id > high {
		id >>= 1
	}
	node, _ := snowflake.NewNode(id)
	return node
}

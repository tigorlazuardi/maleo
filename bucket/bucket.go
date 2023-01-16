package bucket

import (
	"context"
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

// --8<-- [start:bucket]

type Bucket interface {
	// Upload File(s) to the bucket.
	//
	// Implementor must check if File.Data() implements io.Closer and call Close() on it after the upload is complete.
	//
	// Whether the upload is successful or not, the implementor must call File.Close() on all files received.
	//
	// UploadResult.File must be filled with File of the same index and UploadResult.Error must be filled with error if
	// upload operation fails.
	//
	// The number of UploadResult must be the same as the number of files in the parameter.
	Upload(ctx context.Context, files []File) []UploadResult
}

// --8<-- [end:bucket]

package proto

import (
	"os"
	"strings"
	"testing"
)

func TestMerge(t *testing.T) {
	wd, _ := os.Getwd()
	var apiFilePath = strings.SplitN(wd, "/internal", 2)[0] + "/example/proto/1_normal.proto"
	var basePath = strings.SplitN(wd, "/internal", 2)[0] + "/example/proto/base.proto"
	MergeFile(basePath, apiFilePath)
}

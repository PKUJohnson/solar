package fileutil

import (
	"fmt"
	"testing"
)

func TestToken(t *testing.T) {

	//t := GetFileToken("bucket", "key", 123, 456)

	fi, err := FileInfoFromToken("ZpxCLfUP3Rk7FAWRmTasA+0uRub3F2MuxsyRBMRw9ydGY+f434/yy+G3ipbsL07vJAZGlR63B/dnQ+DW9mYUhZQOqfJhdYUGfhRkpTIA5/ljYQPDxig/8KOutZH7rz8POW3NgODEzFj0gKgVgLlW2rbbco9jiS0iydzOJb1b3YexS9kvu098YPPs0rwGSHsRAki0FzrWqrOsfoOn6LZqDg==")
	fmt.Println(fi.ViewerId, fi.Key, fi.Bucket)
	fmt.Println(err)
}

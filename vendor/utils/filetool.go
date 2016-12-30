package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func md5file(filepath string) {
	file, inerr := os.Open(filepath)
	if inerr == nil {
		md5h := md5.New()
		io.Copy(md5h, file)
		fmt.Printf("%x", md5h.Sum([]byte(""))) //md5
	}
}

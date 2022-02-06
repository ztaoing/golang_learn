package main

import (
	"compress/gzip"
	"log"
	"os"
)

// 这个例子中使用gzip压缩格式，标准库还支持zlib, bz2, flate, lzw

func main() {
	outputFile, err := os.Create("test.txt.gz")
	if err != nil {
		log.Fatal(err)
	}

	gzipWriter := gzip.NewWriter(outputFile)
	defer gzipWriter.Close()

	// 当我们写如到gizp writer数据时，它会依次压缩数据并写入到底层的文件中。
	// 我们不必关心它是如何压缩的，还是像普通的writer一样操作即可。
	_, err = gzipWriter.Write([]byte("Gophers rule!\n"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Compressed data written to file.")
}

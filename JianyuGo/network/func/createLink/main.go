package main

import (
	"fmt"
	"log"
	"os"
)

/**
[硬链接]
一个普通的文件是一个指向硬盘的inode的地方。硬链接创建一个新的指针指向同一个地方。
只有所有的链接被删除后文件才会被删除。硬链接只在相同的文件系统中才工作。
你可以认为一个硬链接是一个正常的链接。

[软连接]
symbolic link，又叫软连接，和硬链接有点不一样，它不直接指向硬盘中的相同的地方，而是通过名字引用其它文件。
他们可以指向不同的文件系统中的不同文件。并不是所有的操作系统都支持软链接。
*/
func main() {
	// 创建一个硬链接。
	// 创建后,同一个文件内容会有两个文件名，改变一个文件的内容,会影响另一个。
	// 删除和重命名不会影响另一个。
	err := os.Link("test.txt", "original_also.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("creating sym")
	// Create algorithm symlink
	err = os.Symlink("original.txt", "original_sym.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Lstat返回一个文件的信息，但是当文件是一个软链接时，它返回软链接的信息，而不是引用的文件的信息。
	// Symlink在Windows中不工作。
	fileInfo, err := os.Lstat("original_sym.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Link info: %+v", fileInfo)

	//改变软链接的拥有者不会影响原始文件。
	err = os.Lchown("original_sym.txt", os.Getuid(), os.Getgid())
	if err != nil {
		log.Fatal(err)
	}

	fileInfoS, err := os.Lstat("original_sym.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("new Link info: %+v", fileInfoS)
}

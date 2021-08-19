/**
* @Author:zhoutao
* @Date:2021/5/14 下午1:16
* @Desc:将数据分为header和body
 */

package socket

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

//encode将消息编码
func Encode(message string) ([]byte, error) {
	//length是32位，4个字节
	var length = int32(len(message))
	var pkg = new(bytes.Buffer)
	//写入消息头
	err := binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	//写入body
	err = binary.Write(pkg, binary.LittleEndian, []byte(message))
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), err
}

//decode解码消息
func Decode(reader *bufio.Reader) (string, error) {
	//读取消息的长度,读取前4个字节的数据
	lengthByte, _ := reader.Peek(4)
	//lengthByte为body的长度
	lengthBuff := bytes.NewBuffer(lengthByte)
	var length int32
	//将lengthBuff中的数据读出，存储到length中
	err := binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		return "", err
	}
	//已经被缓存在buffer中的数据的字节数
	if int32(reader.Buffered()) < length+4 {
		return "", err
	}
	//读取真正的数据,header的长度+body的长度
	pack := make([]byte, int(4+length))
	_, err = reader.Read(pack)
	if err != nil {
		return "", err
	}
	return string(pack[4:]), nil
}

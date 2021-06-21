//
// Created by Rustle Karl on 2021.01.05 16:34.
//

package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type People struct {
	Name string
	Age  int
}

func main() {

	p := People{
		Name: "Jason Yin",
		Age:  18,
	}

	/**
	  定义一个字节容器,其结构体如下所示:
	      type Buffer struct {
	          buf      []byte // contents are the bytes buf[off : len(buf)]
	          off      int    // read at &buf[off], write at &buf[len(buf)]
	          lastRead readOp // last read operation, so that Unread* can work correctly.
	      }
	*/
	buf := bytes.Buffer{}

	/**
	  初始化编码器,其函数签名如下:
	      func NewEncoder(w io.Writer) *Encoder
	  以下是函数签名的参数说明:
	      w:
	          一个io.Writer对象,我们可以传递"bytes.Buffer{}"的引用地址
	      返回值:
	          NewEncoder返回将在io.Writer上传输的新编码器。
	*/
	encoder := gob.NewEncoder(&buf)

	/**
	  编码操作
	*/
	err := encoder.Encode(p)
	if err != nil {
		fmt.Println("编码失败,错误原因: ", err)
		return
	}

	/**
	  查看编码后的数据(gob序列化其实是二进制数据哟~)
	*/
	fmt.Println(string(buf.Bytes()))
}

type Student struct {
	Name string
	Age  int
}

func main2() {

	s1 := Student{
		Name: "Jason Yin",
		Age:  18,
	}

	buf := bytes.Buffer{}

	/**
	  初始化编码器
	*/
	encoder := gob.NewEncoder(&buf)

	/**
	  编码操作,相当于序列化过程哟~
	*/
	err := encoder.Encode(s1)
	if err != nil {
		fmt.Println("编码失败,错误原因: ", err)
		return
	}

	/**
	  查看编码后的数据(gob序列化其实是二进制数据哟~)
	*/
	//fmt.Println(string(buf.Bytes()))

	/**
	  初始化解码器,其函数签名如下:
	      func NewDecoder(r io.Reader) *Decoder
	  以下是函数签名的参数说明:
	      r：
	          一个io.Reader对象,其函数签名如下所示:
	              func NewReader(b []byte) *Reader
	          综上所述，我们可以将编码后的字节数组传递给该解码器
	      返回值:
	          new decoder返回从io.Reader读取的新解码器。
	          如果r不同时实现io.ByteReader，它将被包装在bufio.Reader中。
	*/
	decoder := gob.NewDecoder(bytes.NewReader(buf.Bytes()))

	var s2 Student
	fmt.Println("解码之前s2 = ", s2)

	/**
	  进行解码操作，相当于反序列化过程哟~
	*/
	decoder.Decode(&s2)
	fmt.Println("解码之前s2 = ", s2)
}

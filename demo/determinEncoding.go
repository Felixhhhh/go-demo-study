package main

//
//import (
//	"bufio"
//	"fmt"
//	"golang.org/x/net/html/charset"
//	"golang.org/x/text/encoding"
//	"golang.org/x/text/encoding/unicode"
//	"golang.org/x/text/transform"
//	"io/ioutil"
//	"log"
//	"net/http"
//)
//
////处理获取的数据
//func determiEncoding(r *bufio.Reader) encoding.Encoding { //Encoding编码是一种字符集编码，可以在 UTF-8 和 UTF-8 之间进行转换
//	//获取数据,Peek返回输入流的下n个字节
//	bytes, err := r.Peek(1024)
//	if err != nil {
//		log.Printf("fetch error :%v", err)
//		return unicode.UTF8
//	}
//	//调用DEtermineEncoding函数，确定编码通过检查最多前 1024 个字节的内容和声明的内容类型来确定 HTML 文档的编码。
//	e, _, _ := charset.DetermineEncoding(bytes, "")
//	return e
//}
//
//func main() {
//	//res 为结构体，储存了很多的信息
//	resp, err := http.Get("https://www.toutiao.com/?wid=1628221487217")
//	if err != nil {
//		fmt.Println(err)
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//		fmt.Printf("Error status Code :%d", resp.StatusCode)
//	}
//
//	//获取响应体
//	bodyReader := bufio.NewReader(resp.Body)
//
//	//使用determiEncoding函数对获取的信息进行解析
//	e := determiEncoding(bodyReader)
//	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
//
//	//读取并打印获取的信息
//	result, err := ioutil.ReadAll(utf8Reader)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf("%s", result)
//
//}

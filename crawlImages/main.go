package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)


func main() {
	//指定爬取起始、终止页
	var start, end int
	fmt.Println("请输入爬取的起始页(>=1):")
	fmt.Scan(&start)
	fmt.Println("请输入爬取的终止页(>=start):")
	fmt.Scan(&end)

	working(start, end)
}

//获取一个网页所有的内容
func HttpGet(url string) (result string, err error) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		err = err1
		return
	}
	defer resp.Body.Close()
	buf := make([]byte, 4096)
	for {
		n, err2 := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		if err2 != nil && err2 != io.EOF {
			err = err2
			return
		}
		result += string(buf[:n])

	}
	return
}

func working(start, end int) {
	fmt.Printf("正在爬取 %d 到 %d \n", start, end)

	page := make(chan int) //设置多线程

	for i := start; i <= end; i++ {
		go SpidePage(i, page)
	}
	for i := start; i <= end; i++ {
		fmt.Printf("第 %d 页爬取完毕\n", <-page)
	}
}

func SpidePage(i int, page chan int) {
	//网站每一页的改变
	url := "https://wall.alphacoders.com/by_category.php?id=33&name=Women+Wallpapers" + strconv.Itoa(i*1)
	//读取这个页面的所有信息
	result, err := HttpGet(url)
	//判断是否出错，并打印信息
	if err != nil {
		fmt.Println("SpidePage err:", err)
	}

	//正则表达式提取信息
	str := "<div \n    class=\"thumb-container-big \n        \" \n    id=\"thumb_(.*)?\">"
	//解析、编译正则
	ret := regexp.MustCompile(str)
	//提取需要信息-每一个图片的数字
	urls := ret.FindAllStringSubmatch(result, -1)

	for _, jokeURL := range urls {
		//组合每个图片的url
		joke := `https://wall.alphacoders.com/big.php?i=` + jokeURL[1]

		//爬取图片的url
		tuUrl, err := SpideJokePage(joke)
		if err != nil {
			fmt.Println("tuUrl err:", err)
			continue
		}
		SaveJokeFile(tuUrl)
	}

	//防止主go程提前结束
	page <- i
}

//写入文件
func SaveJokeFile(url string) {
	//保存图片的名字是随机的6个数字
	rand.Seed(time.Now().UnixNano())
	var captcha string
	for i := 0; i < 6; i++ {
		//产生0到9的整数
		num := rand.Intn(10)
		//将整数转为字符串
		captcha += strconv.Itoa(num)
	}

	//保存在本地的地址
	path := "imgs/" + captcha + "张.jpg"
	f, err := os.Create(path)
	if err != nil {
		fmt.Println("HttpGet err:", err)
		return
	}
	defer f.Close()

	//读取url的信息
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("http err:", err)
		return
	}
	defer f.Close()

	buf := make([]byte, 4096)
	for {
		n, err2 := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		if err2 != nil && err2 != io.EOF {
			err = err2
			return
		}
		//写入文件
		f.Write(buf[:n])
	}
	fmt.Println("图片: ", path, "下载成功")
}

//爬取图片放大的页面
func SpideJokePage(url string) (tuUrl string, err error) {
	//爬取网站的信息
	result, err1 := HttpGet(url)
	if err1 != nil {
		err = err1
		fmt.Println("SpidePage err:", err)
	}

	str := `<img class="main-content" width="(?s:(.*?))" height="(?s:(.*?))" src="(?s:(.*?))"`
	//解析、编译正则
	ret := regexp.MustCompile(str)
	//提取需要信息-每一个段子的url
	alls := ret.FindAllStringSubmatch(result, -1)
	for a, temTitle := range alls {
		print(a)
		tuUrl = temTitle[3]
		break
	}

	return
}


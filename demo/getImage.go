package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
	_ "io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var headers = map[string][]string{
	"Accept":                    []string{"text/html,application/xhtml+xml,application/xml", "q=0.9,image/webp,*/*;q=0.8"},
	"Accept-Encoding":           []string{"gzip, deflate, sdch"},
	"Accept-Language":           []string{"zh-CN,zh;q=0.8,en;q=0.6,zh-TW;q=0.4"},
	"Accept-Charset":            []string{"utf-8"},
	"Connection":                []string{"keep-alive"},
	"DNT":                       []string{"1"},
	"Host":                      []string{"www.kongjie.com"},
	"Referer":                   []string{"http://www.kongjie.com/home.php?mod=space&do=album&view=all&order=hot&page=1"},
	"Upgrade-Insecure-Requests": []string{"1"},
	"User-Agent":                []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"},
}

func getReponseWithGlobalHeaders(url string) *http.Response {
	req, _ := http.NewRequest("GET", url, nil)
	if headers != nil && len(headers) != 0 {
		for k, v := range headers {
			for _, val := range v {
				req.Header.Add(k, val)
			}
		}
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Sprint(err)
	}
	return res
}

func getHtmlFromUrl(url string) []byte {
	response := getReponseWithGlobalHeaders(url)

	reader := response.Body
	// 返回的内容被压缩成gzip格式了，需要解压一下
	if response.Header.Get("Content-Encoding") == "gzip" {
		reader, _ = gzip.NewReader(response.Body)
	}
	// 此时htmlContent还是gbk编码，需要转换成utf8编码
	htmlContent, _ := ioutil.ReadAll(reader)

	oldReader := bufio.NewReader(bytes.NewReader(htmlContent))
	peekBytes, _ := oldReader.Peek(1024)
	e, _, _ := charset.DetermineEncoding(peekBytes, "")
	utf8reader := transform.NewReader(oldReader, e.NewDecoder())
	// 此时htmlContent就已经是utf8编码了
	htmlContent, _ = ioutil.ReadAll(utf8reader)

	if err := response.Body.Close(); err != nil {
		fmt.Println("error happened when closing response body!", err)
	}
	return htmlContent
}

func HandleError(err error, why string) {
	if err != nil {
		fmt.Println(why, err)
	}
}

// 下载图片，传入的是图片叫什么
func DownloadFile(url string, filename string) (ok bool) {
	resp, err := http.Get(url)
	result := getHtmlFromUrl(url)
	fmt.Println(result)
	HandleError(err, "http.get.url")
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)

	HandleError(err, "resp.body")
	filename = "img/" + filename
	// 写出数据
	err = ioutil.WriteFile(filename, bytes, 0666)
	if err != nil {
		return false
	} else {
		return true
	}
}

// 并发爬思路：
// 1.初始化数据管道
// 2.爬虫写出：26个协程向管道中添加图片链接
// 3.任务统计协程：检查26个任务是否都完成，完成则关闭数据管道
// 4.下载协程：从管道里读取链接并下载

var (
	// 存放图片链接的数据管道
	chanImageUrls chan string
	waitGroup     sync.WaitGroup
	// 用于监控协程
	chanTask chan string
	reImg    = `https?://[^"]+?(\.((jpg)|(png)|(jpeg)|(gif)|(bmp)))`
)

func main() {
	// DownloadFile("http://i1.shaodiyejin.com/uploads/tu/201909/10242/e5794daf58_4.jpg", "1.jpg")

	// 1.初始化管道
	chanImageUrls = make(chan string, 1000000)
	chanTask = make(chan string, 26)
	// 2.爬虫协程
	for i := 1; i < 27; i++ {
		waitGroup.Add(1)
		go getImgUrls("https://www.bizhizu.cn/shouji/tag-%E5%8F%AF%E7%88%B1/" + strconv.Itoa(i) + ".html")
	}
	// 3.任务统计协程，统计26个任务是否都完成，完成则关闭管道
	waitGroup.Add(1)
	go CheckOK()
	// 4.下载协程：从管道中读取链接并下载
	for i := 0; i < 5; i++ {
		waitGroup.Add(1)
		go DownloadImg()
	}
	waitGroup.Wait()
}

// 下载图片
func DownloadImg() {
	for url := range chanImageUrls {
		filename := GetFilenameFromUrl(url)
		ok := DownloadFile(url, filename)
		if ok {
			fmt.Printf("%s 下载成功\n", filename)
		} else {
			fmt.Printf("%s 下载失败\n", filename)
		}
	}
	waitGroup.Done()
}

// 截取url名字
func GetFilenameFromUrl(url string) (filename string) {
	// 返回最后一个/的位置
	lastIndex := strings.LastIndex(url, "/")
	// 切出来
	filename = url[lastIndex+1:]
	// 时间戳解决重名
	timePrefix := strconv.Itoa(int(time.Now().UnixNano()))
	filename = timePrefix + "_" + filename
	return
}

// 任务统计协程
func CheckOK() {
	var count int
	for {
		url := <-chanTask
		fmt.Printf("%s 完成了爬取任务\n", url)
		count++
		if count == 26 {
			close(chanImageUrls)
			break
		}
	}
	waitGroup.Done()
}

// 爬图片链接到管道
// url是传的整页链接
func getImgUrls(url string) {
	urls := getImgs(url)
	// 遍历切片里所有链接，存入数据管道
	for _, url := range urls {
		chanImageUrls <- url
	}
	// 标识当前协程完成
	// 每完成一个任务，写一条数据
	// 用于监控协程知道已经完成了几个任务
	chanTask <- url
	waitGroup.Done()
}

// 获取当前页图片链接
func getImgs(url string) (urls []string) {
	pageStr := GetPageStr(url)
	re := regexp.MustCompile(reImg)
	results := re.FindAllStringSubmatch(pageStr, -1)
	fmt.Printf("共找到%d条结果\n", len(results))
	for _, result := range results {
		url := result[0]
		urls = append(urls, url)
	}
	return
}

// 抽取根据url获取内容
func GetPageStr(url string) (pageStr string) {
	resp, err := http.Get(url)
	HandleError(err, "http.Get url")
	defer resp.Body.Close()
	// 2.读取页面内容
	pageBytes, err := ioutil.ReadAll(resp.Body)
	HandleError(err, "ioutil.ReadAll")
	// 字节转字符串
	pageStr = string(pageBytes)
	return pageStr
}

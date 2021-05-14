/**
* @Author:zhoutao
* @Date:2021/5/14 上午8:56
* @Desc:断点下载
 */

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

//filePart 文件分片
type filePart struct {
	Index int    //文件分片的序号
	From  int    // 开始的byte
	To    int    //结束的byte
	Data  []byte //http下载得到的文件内容
}

//文件下载器
type FileDownloader struct {
	fileSize       int
	url            string
	outputFileName string
	outputDir      string
	totalPart      int        //下载的线程数
	doneFilePart   []filePart //存储了已经完成的分片

}

//创建一个request
func (d FileDownloader) getNewRequest(method string) (*http.Request, error) {
	req, err := http.NewRequest(method, d.url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "tao")
	return req, nil
}

// 下载分片
func (d FileDownloader) downloadPart(c filePart) error {
	r, err := d.getNewRequest("GET")
	if err != nil {
		return err
	}
	//开始下载指定范围的分片
	log.Printf("开始下载[%d]下载from：%d to:%d\n", c.Index, c.From, c.To)
	r.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", c.From, c.To))
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	if resp.StatusCode > 299 {
		return errors.New(fmt.Sprintf("服务器状态码：%v", resp.StatusCode))
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(bs) != (c.To - c.From + 1) {
		return errors.New("下载文件分片长度错误")
	}
	c.Data = bs
	//加入到已接收中
	d.doneFilePart[c.Index] = c
	return nil
}

// 合并下载的文件
func (d FileDownloader) mergeFileParts() error {
	log.Println("开始合并文件")
	path := filepath.Join(d.outputDir, d.outputFileName)
	mergedFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer mergedFile.Close()
	hash := sha256.New()
	totalSize := 0
	for _, s := range d.doneFilePart {
		//写入文件
		mergedFile.Write(s.Data)
		//
		hash.Write(s.Data)
		totalSize += len(s.Data)
	}
	if totalSize != d.fileSize {
		return errors.New("文件不完整")
	}
	//校验 sha256校验
	if hex.EncodeToString(hash.Sum(nil)) != "" {
		return errors.New("文件损坏")
	} else {
		log.Println("文件SHA-256校验成功")
	}
	return nil
}

func NewFileDownloader(url string, outputFileName, outputDir string, totalPart int) *FileDownloader {
	if outputDir == "" {
		//获取当前工作目录
		wd, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}
		outputDir = wd
	}
	return &FileDownloader{
		fileSize:       0,
		url:            url,
		outputFileName: outputFileName,
		totalPart:      totalPart,
		doneFilePart:   make([]filePart, totalPart),
		outputDir:      outputDir,
	}
}

func parseFileinfo(resp *http.Response) string {
	contentDisposition := resp.Header.Get("Content-Disposition")
	if contentDisposition != "" {
		_, params, err := mime.ParseMediaType(contentDisposition)
		if err != nil {
			panic(err)
		}
		return params["filename"]
	}
	// build filename
	filename := filepath.Base(resp.Request.URL.Path)
	return filename
}

// 获取要下载的文件的基本信息，使用http method head
func (d *FileDownloader) head() (int, error) {
	r, err := d.getNewRequest("HEAD")
	if err != nil {
		return 0, err
	}
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode > 299 {
		return 0, errors.New(fmt.Sprintf("can't process ,response is %v", resp.StatusCode))
	}
	//减产是否支持断点续传
	if resp.Header.Get("Accept-Ranges") != "bytes" {
		return 0, errors.New("服务器不支持断点续传")
	}
	//https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Ranges
	d.outputFileName = parseFileinfo(resp)
	//https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Length
	return strconv.Atoi(resp.Header.Get("Content-Length"))
}

func (d *FileDownloader) Run() error {
	fileTotalSize, err := d.head()
	if err != nil {
		return err
	}
	d.fileSize = fileTotalSize
	jobs := make([]filePart, d.totalPart)
	//每个分片的大小
	eachSize := fileTotalSize / d.totalPart
	//构造每一个jobs
	for i := range jobs {
		jobs[i].Index = i
		//定义from 开始
		if i == 0 {
			jobs[i].From = 0
		} else {
			jobs[i].From = jobs[i-1].To + 1
		}
		//定义to 结尾
		if i < d.totalPart-1 {
			jobs[i].To = jobs[i].From + eachSize
		} else {
			//最后一个分片
			jobs[i].To = fileTotalSize - 1
		}
	}

	var wg sync.WaitGroup
	for _, j := range jobs {
		wg.Add(1)
		go func(job filePart) {
			defer wg.Done()
			err := d.downloadPart(job)
			if err != nil {
				log.Println("下载文件失败", err, job)
				//TODO将下载失败的job加入failed job中进行重试
			}
		}(j)
	}
	//TODO 如果failed job不为空，则需要处理下载失败的job
	wg.Wait()
	//合并
	return d.mergeFileParts()
}

func main() {
	startTime := time.Now()
	var url string
	url = "https://download.jetbrains.com.cn/go/goland-2021.1.1.dmg"
	downloader := NewFileDownloader(url, "", "", 10)
	if err := downloader.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n 文件下载完成耗时：%f second\n", time.Now().Sub(startTime).Seconds())
}

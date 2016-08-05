package main

import (
	"os"
	"github.com/smtc/glog"
	"bufio"
	"io"
	"strings"
	"net/url"
	"net/http"
	"io/ioutil"
)

/**
打开文件并处理
创建人:邵炜
创建时间:2016年6月1日09:40:03
输入参数:filePath 文件地址
输出参数:文件对象 错误对象
*/
func openFile(filePath string) (*os.File, error) {
	var (
		fs  *os.File
		err error
	)
	fs, err = os.Open(filePath)
	if err != nil {
		glog.Error("open file is error! filePath: %s err: %s \n", filePath, err.Error())
		return nil, err
	}
	glog.Info("file open success! filePath: %s \n", filePath)
	return fs, nil
}

/**
文件读取方法
创建人:邵炜
创建时间:2016年8月5日11:25:08
输入参数:文件路劲
 */
func readFileMobs(path string) {
	var (
		readAll     = false
		readByte    []byte
		line        []byte
		err         error
	)
	read, err := openFile(path)
	if err != nil {
		return
	}
	defer read.Close()
	buf := bufio.NewReader(read)
	for err != io.EOF {
		if err != nil {
			glog.Error("read error! err: %s \n", err.Error())
		}
		if readAll {
			readByte, readAll, err = buf.ReadLine()
			line = append(line, readByte...)
		} else {
			readByte, readAll, err = buf.ReadLine()
			line = append(line, readByte...)
			if len(strings.TrimSpace(string(line))) == 0 {
				continue
			}
			go mobOrderInterFace(string(line))
			line = line[:0]
		}
	}
}

/**
订购接口调用
创建人:邵炜
创建时间:2016年8月5日11:28:02
输入参数:输入参数内容
 */
func mobOrderInterFace(content string) {
urlStr:=interFaceApi+"/DailyOrderProduct"
	contentArray:=strings.Split(content,",")
	valueData:=url.Values{}
	valueData.Add("mob",contentArray[0])
	valueData.Add("operType","0")
	valueData.Add("productId",contentArray[1])
	valueData.Add("effectMode","1")
	valueData.Add("serviceType","0000")
	valueData.Add("changeType","1")
	valueData.Add("packageId",contentArray[2])
	valueData.Add("elementId",contentArray[3])
	valueData.Add("elementType","D")
	valueData.Add("password","njaction@2012")
	resp,err:=http.PostForm(urlStr,valueData)
	if err != nil {
		glog.Error("mobOrderInterFace http post error! urlStr: %s valueData: %v err: %s \n",urlStr,valueData,err.Error())
		return
	}
	bodyContent,err:=ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		glog.Error("mobOrderInterFace http post can't read ! urlStr: %s valueData: %v err: %s \n",urlStr,valueData,err.Error())
		return
	}
	glog.Info("mobOrderInterFace http success! urlStr: %s valueData: %v content: %s \n",urlStr,valueData,string(bodyContent))
}
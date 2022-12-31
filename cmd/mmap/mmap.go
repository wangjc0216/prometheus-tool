package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/prometheus/tsdb/fileutil"
	"github.com/wangjc/prometheus-tool/pkg/logger"
	"net/http"
	"net/url"
	"strconv"
)

func init() {
	logger.Init("mmap-webserver.log", "info", "mmap-webserver", true)
}

var mmapFile *fileutil.MmapFile
var mmapBytes []byte

func main() {
	filePath := flag.String("filepath", "./mmapfile.txt", "the filepath of mmap file")
	flag.Parse()
	var err error
	mmapFile, err = fileutil.OpenMmapFile(*filePath)
	if err != nil {
		logger.Errorf("oepnMmapFile error:%v", err)
		return
	}
	//这里的bs虽然看起来很大(len cap都是5000000000)，但是并没有全部载入到物理内存中
	mmapBytes = mmapFile.Bytes()
	fmt.Println("first 10 bytes is %v", mmapBytes[:10])

	http.HandleFunc("/mmap", mmapHandle)
	go func() {
		http.ListenAndServe(":9001", promhttp.Handler())
	}()
	http.ListenAndServe(":8001", nil)

}

func mmapHandle(resp http.ResponseWriter, req *http.Request) {
	values, _ := url.ParseQuery(req.URL.RawQuery)
	sizeStr := values.Get("size")
	offsetStr := values.Get("offset")

	var size, offset int
	var err error
	if sizeStr == "" {
		size = 10000
	} else {
		size, err = strconv.Atoi(sizeStr)
		if err != nil {
			logger.Errorf("sizeStr(%v) cannot convert to number,error is %v", sizeStr, err)
			size = 10000
		}
	}

	if offsetStr == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			logger.Errorf("offsetStr(%v) cannot convert to number,error is %v", offsetStr, err)
			offset = 0
		}
	}
	logger.Infof("offset is %v,size is %v", offset, size)
	bs := mmapBytes[offset : offset+size]
	var count int
	for _, b := range bs {
		if b == 0 {
			count++
		}
	}
	logger.Infof("byte 0 counter is %v ", count)
}

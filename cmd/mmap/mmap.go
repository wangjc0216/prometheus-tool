package main

import (
	"fmt"
	"github.com/prometheus/prometheus/tsdb/fileutil"
	"github.com/wangjc/prometheus-tool/pkg/logger"
	"time"
)

func main() {
	path := "./mmapfile.txt"
	f, err := fileutil.OpenMmapFile(path)
	if err != nil {
		logger.Errorf("oepnMmapFile error:%v", err)
		return
	}
	//这里的bs虽然看起来很大(len cap都是5000000000)，但是并没有全部载入到物理内存中
	bs := f.Bytes()
	fmt.Println("first 10 bytes is %v", bs[:10])
	time.Sleep(time.Minute * 30)

}

package main

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"github.com/prometheus/prometheus/tsdb/chunks"
	"github.com/prometheus/prometheus/tsdb/record"
	"github.com/prometheus/prometheus/tsdb/wal"
	"github.com/wangjc/prometheus-tool/pkg/logger"
)

func logInit(printFlag bool) {
	logger.Init("wal-tool.log", "info", "wal", printFlag)
}

type metric struct {
	Series  record.RefSeries
	Sampels []record.RefSample
}

func descSegmentFile(segFile string) error {
	if segFile == "" {
		return nil
	}
	seg, err := wal.OpenReadSegment(segFile)
	if err != nil {
		return errors.Wrap(err, "OpenReadSegment error")
	}
	reader := wal.NewReader(seg)
	var (
		dec       record.Decoder
		series    []record.RefSeries
		samples   []record.RefSample
		exemplars []record.RefExemplar
	)

	for reader.Next() {
		rec := reader.Record()

		switch dec.Type(rec) {
		case record.Series:
			var newSeries []record.RefSeries
			newSeries, err = dec.Series(rec, newSeries)
			series = append(series, newSeries...)

		case record.Samples:
			var newSamples []record.RefSample
			newSamples, err = dec.Samples(rec, newSamples)
			samples = append(samples, newSamples...)

		case record.Exemplars:
			var newExamplars []record.RefExemplar
			newExamplars, err = dec.Exemplars(rec, newExamplars)
			exemplars = append(exemplars, newExamplars...)

		case record.Tombstones:
			logger.Infof("tombstones...")

		case record.Unknown:
			logger.Infof("unknown...")
		default:
		}
		if err != nil {
			return errors.Wrap(err, "decode error")
		}
	}

	metricInfo := make(map[chunks.HeadSeriesRef]*metric)

	//Return error if found there are duplicated series
	for _, s := range series {
		if _, ok := metricInfo[s.Ref]; ok {
			return errors.Wrapf(err, "metric Series exist,series is %+v", s)
		}
		metricInfo[s.Ref] = &metric{
			Series: s,
		}
	}
	//Return error if found samples is not in range of series
	for _, sample := range samples {
		if _, ok := metricInfo[sample.Ref]; !ok {
			//return errors.Wrapf(errors.New("sample doesn't exist"), "metric Sample doesn't exist, sample is %+v", sample)

			logger.Errorf("sample doesn't exist,continue")
			continue
		}
		// if sampels is nil , alloc
		if metricInfo[sample.Ref].Sampels == nil {
			metricInfo[sample.Ref].Sampels = make([]record.RefSample, 0)
		}
		metricInfo[sample.Ref].Sampels = append(metricInfo[sample.Ref].Sampels, sample)
	}

	for _, m := range metricInfo {
		//todo outputformat can be better
		fmt.Println(m)
	}
	return nil
}
func descDir(walDir string) error {
	fmt.Printf("walDir is %s\n", walDir)
	f, l, err := wal.Segments(walDir)
	if err != nil {
		return errors.Wrap(err, "wal segments error")
	}

	for i := f; i <= l; i++ {
		segmentName := wal.SegmentName(walDir, i)
		fmt.Println(segmentName)
	}
	logger.Infof("first segmentNo is %d  and last segmentNo is %d", f, l)
	return nil
}

func main() {
	defaultDir := "./mock_data/wal"

	walDir := flag.String("wal_dir", defaultDir, "the directory of wal file")
	segFile := flag.String("segfile", "", "the segment file of wal directory")
	print := flag.Bool("print", false, "whether print log to console")
	flag.Parse()
	logInit(*print)

	err := descDir(*walDir)
	if err != nil {
		//fix:
		fmt.Println("descDir error:%v", err)   //用于输出终端
		logger.Errorf("descDir error:%v", err) //用于输出到日志
		return
	}

	err = descSegmentFile(*segFile)
	if err != nil {
		fmt.Println("descSegmentFile error:%v", err)   //用于输出终端
		logger.Errorf("descSegmentFile error:%v", err) //用于输出到日志
		return
	}

}

//todo some question
//1. RefId 是如何定义的
//2. replay的时候需要注意什么呢
//3. wal的数据是如何做到有多种数据类型的？咋会将series和samples杂糅在一个文件里呢？
//4. 对于终端命令行相关的工具，我需要日志是打印在终端的是可以简洁的，像fmt.print的；打印在文件的是详细的
//5. 同一个segment，会有没有对应series信息的sample，这个series和sample是如何处理的呢？我的理解是series类似于metadata，它是多久保存一次呢？

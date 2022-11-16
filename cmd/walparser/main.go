package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pkg/errors"
	"github.com/prometheus/prometheus/tsdb"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
	"github.com/prometheus/prometheus/tsdb/chunks"
	"github.com/prometheus/prometheus/tsdb/index"
	"github.com/prometheus/prometheus/tsdb/record"
	"github.com/prometheus/prometheus/tsdb/wal"
	"github.com/wangjc/prometheus-tool/pkg/logger"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func logInit(printFlag bool) {
	logger.Init("wal-tool.log", "info", "wal", printFlag)
}

type metric struct {
	Series  record.RefSeries
	Samples []record.RefSample
}

type metricCollection map[chunks.HeadSeriesRef]*metric

func (mc metricCollection) tableOuput() {
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"HeadSeriesRef", "Labels", "Timestamp", "Value"})
	count := 0
	for ref, m := range mc {
		for _, sample := range m.Samples {
			t.AppendRow(table.Row{ref, m.Series.Labels, sample.T, sample.V}, rowConfigAutoMerge)
		}
		t.AppendSeparator()
		//just print top 50 of the metric
		if count > 50 {
			t.AppendFooter(table.Row{fmt.Sprintf("just print top 50,total number of series is %d", len(mc))})
			break
		} else {
			count++
		}
	}
	t.Render()
}

//描述SegmentFile
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

	var metricInfo metricCollection = make(map[chunks.HeadSeriesRef]*metric)

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
		if metricInfo[sample.Ref].Samples == nil {
			metricInfo[sample.Ref].Samples = make([]record.RefSample, 0)
		}
		metricInfo[sample.Ref].Samples = append(metricInfo[sample.Ref].Samples, sample)
	}

	metricInfo.tableOuput()

	return nil
}

//描述Wal directory
func descDir(walDir string) error {
	if walDir == "" {
		return nil
	}
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

//描述block index file
func descBlockFile(blockPath string) error {
	if blockPath == "" {
		return nil
	}
	indexPath := filepath.Join(blockPath, "index")
	indexer, err := index.NewFileReader(indexPath)
	if err != nil {
		return err
	}
	//todo 1需要了解index 文件格式是怎样的，都有哪些组成的，如何反序列化出来
	//todo 2当mmap的时候，是会把index文件所有数据load到内存中吗？应该不是这样，因为我看到chunkReader 也是MMapOpen的
	//todo 3 index.Reader接口 看了还有其他的实现，其他的实现是在哪里用到的呢？

	names, err := indexer.LabelNames()
	if err != nil {
		return err
	}
	fmt.Println("labelnames is: ", names)

	values, err := indexer.LabelValues("consumerGroup")
	if err != nil {
		return err
	}
	fmt.Println("labelvalues is:", values)

	return nil
}

//描述 block chunk file
func descChunkFile(blockPath string) error {
	if blockPath == "" {
		return nil
	}
	//todo 每个block都有一个chunks目录，chunks目录下，可以进行读取
	pool := chunkenc.NewPool()
	chunksDir := filepath.Join(blockPath, "chunks")
	cr, err := chunks.NewDirReader(chunksDir, pool)
	if err != nil {
		return err
	}
	chunk, err := cr.Chunk(1)
	if err != nil {
		return err
	}
	_ = chunk
	return nil

}

//描述 block file metadata
func descMetadata(blockPath string) error {
	if blockPath == "" {
		return nil
	}
	b, err := ioutil.ReadFile(filepath.Join(blockPath, "meta.json"))
	if err != nil {
		return err
	}
	var m tsdb.BlockMeta

	err = json.Unmarshal(b, &m)
	if err != nil {
		return err
	}
	fmt.Printf("metadata:%+v\n", m)

	return nil
}

func main() {
	defaultDir := "./mock_data/wal"
	defaultBlockFile := "./mock_data/blockfile/01G97CX7EQ4C5G3XR3CWK1V828"

	walDir := flag.String("wal_dir", "", fmt.Sprintf("the directory of wal file, example %s", defaultDir))
	segFile := flag.String("segfile", "", "the segment file of wal directory")
	blockIndexLocation := flag.String("blockindex", "", fmt.Sprintf("the directory of  block index file,example %s", defaultBlockFile))
	blockChunks := flag.String("blockchunk", "", fmt.Sprintf("the directory of  block chunk file,example %s", defaultBlockFile))
	blockMetadata := flag.String("blockmetadata", "", fmt.Sprintf("the directory of block meta data ,example %s", defaultBlockFile))

	//todo chunk

	print := flag.Bool("print", false, "whether print log to console")
	flag.Parse()
	logInit(*print)
	startTime := time.Now()
	defer func() {
		fmt.Printf("execute Time spent %f seconds\n", time.Now().Sub(startTime).Seconds())
	}()

	err := descDir(*walDir)
	checkErr("descDir", err)

	err = descSegmentFile(*segFile)
	checkErr("descSegmentFile", err)

	err = descBlockFile(*blockIndexLocation)
	checkErr("descBlockFile", err)

	err = descChunkFile(*blockChunks)
	checkErr("chunksLocation", err)

	err = descMetadata(*blockMetadata)
	checkErr("blockMetadata", err)

}

func checkErr(funcname string, err error) {
	if err != nil {
		fmt.Printf("%s error:%v\n", funcname, err)    //用于输出终端
		logger.Errorf("%s error:%v\n", funcname, err) //用于输出到日志
		os.Exit(-1)
	}
	return
}

//todo some question
//1. RefId 是如何定义的
//2. replay的时候需要注意什么呢
//3. wal的数据是如何做到有多种数据类型的？咋会将series和samples杂糅在一个文件里呢？
// 通过storage.Appender接口的Append()方法来做

//4. 对于终端命令行相关的工具，我需要日志是打印在终端的是可以简洁的，像fmt.print的；打印在文件的是详细的
//5. 同一个segment，会有没有对应series信息的sample，这个series和sample是如何处理的呢？我的理解是series类似于metadata，它是多久保存一次呢？

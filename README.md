# prometheus-tool
> 一些基于prometheus周边生态的小工具，可帮助你快速理解prometheus、更好使用prometheus

## 1. expr-ast-tree
对promql进行分析，解析出Expr Tree(词法解析器)

## 2. wal/block/snapshot解析
对prometheus文件进行解析。

## 3. prometheus-alertrule
动态在prometheus进行配置告警规则，并进行reload。
如配置告警周期(告警持续多久告警or告警累计几次会进行告警)则需要在Alertmanager上进行配置。这里需要注意：
(初级使用如告警持续多久后发出可以通过配置alertmanager的yaml，告警累计次数需要进行二次开发得出)


## 4. persistant Alertmanager
对告警消息进行更好的持久化，对数据可通过kafka、pg、es等方式进行记录。

## 5. prometheus operator +
prometheus operator 的增强
相较于prometheus operator 不重视的case，可以适当进行增强，如remote write。

## 6. thanos & M3DB
M3DB的部署实现、关键组件的快速上手、对于性能瓶颈的调优

## 7. prometheus-querylog
针对querylog，可以进行统计、分析、记录，并根据查询情况进行指标暴露。


## 8. 性能调优实现


## 9. Prometheus BlockFile Backfill

1）prometheus可能因为数据retention的缘故导致历史数据会清除，我们如果想long-term分析其指标数据，
需要对数据进行长时间保存。对于M3DB来说，它支持long-term的指标保存，但是运维成本略高，可否有一种思路，
将prometheus的数据取出，在大数据平台中读取。

还有一种可能，使用了远程TSDB，那么则可能会出现TSDB故障的问题，在修复后，需要从边缘prometheus backfill会数据。

2）对指标的数据分析。T+1 对BlockFile进行拉取，并通过Flink 或 批处理


## 10. Prometheus ingest sharding
对指标摄取进行分片。当集群比较大，如几百台机器组成的Kubernetes集群，那么对prometheus不仅仅要部署高可用机器（两台）。可能还需要对target进行分片。


## 11. snapshot与checkpoint



## 12. fanout query engine
如果是单纯的remote read，来做告警（如从多个prometheus实例读取数据），则会依赖最后一个prometheus的返回延时。 如果出现了实例宕机的情况，可能也会对查询结果造成影响。
这里可以参考m3query的查询是如何fanout的，是如何做sharding的，这里涉及到新增节点后的data rebalance ,resharding。


## 13. debug query
从目标prometheus来拉取下Blockfile或者wal数据，分别查出来阶段性数据（处理）情况。


## 14. metric debug 
找到对应metric 是从哪个target获取的，因为中间会有一些relabel的操作，影响寻找。




# TODO

## 数据格式理解
- [x] 解析promql的语法，打印ast(Expr)
- [ ] 解析WAL的内容格式，打印
- [ ] 解析Block 的内容格式，打印


## 可观测性相关Debug

- [ ] 定时load 对应pprof，并进行分析比对，处理
- [ ] 一条promql对应的性能消耗计算（性能消耗可以由自己来进行定义）

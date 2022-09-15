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


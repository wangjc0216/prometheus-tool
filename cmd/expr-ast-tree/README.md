# Usage

>  sum(core_infra_api_total_request{cluster_name="product-pro",namespace="core-infra",service_name="pro-test-flow-compare-1",function_name!="all"})by(cluster_name,namespace,service_name,pod) and on (pod) kube_pod_container_info{container="pro-test-flow-compare-1",monitor_cluster="product-pro"}

> go run main.go --promql 'histogram_quantile(0.99, sum by(cluster_name, service_name, namespace, pod, function_name, le) (rate(infra_api_latency_h_bucket{cluster_name="openplatform-pro",function_name!="all"}[2m])))'

得到Expr的AST树：
```
expr
└── Type:*parser.Call

└── Func:{histogram_quantile [scalar vector] 0 vector}

└── PosRange:{0 193}

└── Args:

    └── arg[0]
    │   ├── Type:*parser.NumberLiteral

    │   ├── [*parser.NumberLiteral]:0.99

    └── arg[1]
        └── Type:*parser.AggregateExpr

        └── Op:SUM

        └── Grouping:[cluster_name service_name namespace pod function_name le]

        └── Without:false

        └── PosRange:{25 192}

        └── Expr:

        │   ├── Type:*parser.Call

        │   ├── Func:{rate [matrix] 0 vector}

        │   ├── PosRange:{96 191}

        │   ├── Args:

        │       └── arg[0]
        │           └── Type:*parser.MatrixSelector

        │           └── [*parser.MatrixSelector]:infra_api_latency_h_bucket{cluster_name="openplatform-pro",function_name!="all"}[2m]

        │           └── VectorSelector:
        │               └── Type:*parser.VectorSelector

        │               └── [*parser.VectorSelector]:infra_api_latency_h_bucket{cluster_name="openplatform-pro",function_name!="all"}

        └── Param:
```



# Introduction

This repo is a fork from https://github.com/jaegertracing/jaeger and add extra otel plugin.

Below is the added plugin:
```shell
go get github.com/open-telemetry/opentelemetry-collector-contrib/connector/servicegraphconnector@v0.127.0
go get github.com/open-telemetry/opentelemetry-collector-contrib/exporter/loadbalancingexporter@v0.127.0
go get github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter@v0.127.0
go get github.com/open-telemetry/opentelemetry-collector-contrib/exporter/sumologicexporter@v0.127.0
go get github.com/open-telemetry/opentelemetry-collector-contrib/processor/probabilisticsamplerprocessor@v0.127.0
go get github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor@v0.127.0
go get github.com/open-telemetry/opentelemetry-collector-contrib/receiver/fluentforwardreceiver@v0.127.0
go get github.com/open-telemetry/opentelemetry-collector-contrib/receiver/skywalkingreceiver@v0.127.0
```

## Sync fork step-to-step

1. trigger "Sync fork" button for **main branch** from repo main page
2. checkout the target branch(v2.7.0) from **main branch** at desire commit
3. checkout a new branch from **insight-main** and named it **upgrade-to-v2.7.0**(any words you like)
4. merge the target branch(v2.7.0) to upgrade-to-v2.7.0 branch 
5. create a PR to merge this upgrade to **insight-main** branch


### Notice

you will get conflict and start resolve in step 4, you can fix it by:
1. for `go.mod` file resolve, accept branch(v2.7.0) firstly and use the above command to add back the dependency from our added 
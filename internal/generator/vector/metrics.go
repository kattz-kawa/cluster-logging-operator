package vector

import (
	"golang.org/x/sys/unix"
)

const (
	InternalMetricsSourceName = "internal_metrics"
	PrometheusOutputSinkName  = "prometheus_output"

	AddNodenameToMetricTransformName = "add_nodename_to_metric"
)

var PrometheusExporterAddress string

func init() {
	if fd, err := unix.Socket(unix.AF_INET6, unix.SOCK_STREAM, unix.IPPROTO_IP); err != nil {
		PrometheusExporterAddress = `0.0.0.0`
	} else {
		unix.Close(fd)
		PrometheusExporterAddress = `[::]`
	}
	PrometheusExporterAddress += `:24231`
}

type InternalMetrics struct {
	ID                string
	ScrapeIntervalSec int
}

func (InternalMetrics) Name() string {
	return "internalMetricsTemplate"
}

// #namespace = "collector"
// #scrape_interval_secs = {{.ScrapeIntervalSec}}
func (i InternalMetrics) Template() string {
	return `
{{define "` + i.Name() + `" -}}
[sources.{{.ID}}]
type = "internal_metrics"
{{end}}
`
}

type PrometheusExporter struct {
	ID            string
	Inputs        string
	Address       string
	TlsMinVersion string
	CipherSuites  string
}

func (p PrometheusExporter) Name() string {
	return "PrometheusExporterTemplate"
}

func (p PrometheusExporter) Template() string {
	return `{{define "` + p.Name() + `" -}}
[sinks.{{.ID}}]
type = "prometheus_exporter"
inputs = {{.Inputs}}
address = "{{.Address}}"
default_namespace = "collector"

[sinks.{{.ID}}.tls]
enabled = true
key_file = "/etc/collector/metrics/tls.key"
crt_file = "/etc/collector/metrics/tls.crt"
min_tls_version = "{{.TlsMinVersion}}"
ciphersuites = "{{.CipherSuites}}"
{{end}}`
}

type AddNodenameToMetric struct {
	ID     string
	Inputs string
}

func (a AddNodenameToMetric) Name() string {
	return AddNodenameToMetricTransformName
}

func (a AddNodenameToMetric) Template() string {
	return `{{define "` + a.Name() + `" -}}
[transforms.{{.ID}}]
type = "remap"
inputs = {{.Inputs}}
source = '''
.tags.hostname = get_env_var!("VECTOR_SELF_NODE_NAME")
'''
{{end}}`
}

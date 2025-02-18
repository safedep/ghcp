package init

import "github.com/safedep/dry/obs"

func init() {
	obs.InitPrometheusMetricsProvider(obs.AppServiceName("ghcp"), "")
}

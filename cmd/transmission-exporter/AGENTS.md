# Collectors

- Name a metric `transmission_<collector>_<name>` from the `namespace` and `collectorNamespace` constants; a new byte metric ends in `_bytes`.
- Metric names and labels are public API: the dashboards in `dashboards/` and users' alerts query them. A rename is a breaking change (`feat!:`), with the dashboards updated in the same change.
- A torrent field the RPC client doesn't fetch yet goes on `Torrent` in the root `torrent.go`; its json tag puts it in the `torrent-get` request.
- When its RPC call fails, `Collect` sends `prometheus.NewInvalidMetric`: `/metrics` answers 500 and Prometheus' `up` drops to 0. Don't add a `transmission_up` gauge, which the 500 would hide.
- `Describe` sends every descriptor `Collect` uses, or `/metrics` fails with HTTP 500. `collector_test.go` catches a missing one when the metric is in its expected output.

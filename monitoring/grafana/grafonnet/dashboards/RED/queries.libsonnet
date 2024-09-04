local grafonnet = import 'github.com/grafana/grafonnet/gen/grafonnet-latest/main.libsonnet';
local prometheusQuery = grafonnet.query.prometheus;

{
  local rate(app_label) =
    prometheusQuery.new(
      'Prometheus',
      |||
        sum by (%s)(
            rate(istio_requests_total{
                %s="${sourceApp}",
                reporter="source"
                }[${interval}]
            )
        ) 
      ||| % [app_label, app_label],
    )
    + prometheusQuery.withIntervalFactor(2)
    + prometheusQuery.withLegendFormat('{{%s}}' % app_label),

  local succsessRate(app_label) =
    prometheusQuery.new(
      'Prometheus',
      |||
        sum by (%s) (
          rate(
            istio_requests_total{
              reporter="source",
              response_code!~"5.*",
              %s="${sourceApp}"
            }[$interval]
          )
          )
          /
          sum by (%s) (
            rate(
              istio_requests_total{
                reporter="source",
                %s="${sourceApp}"
              }[$interval]
            )
          )
      ||| % [app_label, app_label, app_label, app_label],
    )
    + prometheusQuery.withIntervalFactor(2)
    + prometheusQuery.withLegendFormat('{{%s}}' % app_label),

  local duration(app_label) =
    [
      prometheusQuery.new(
        'Prometheus',
        |||
          histogram_quantile(
            0.%s,
            sum by (le, %s) (
              rate(
                istio_request_duration_milliseconds_bucket{
                    reporter=~"source",
                    %s="${sourceApp}",
                }
              [$interval])
            )
          )
        ||| % [quantile, app_label, app_label]
      )
      + prometheusQuery.withIntervalFactor(2)
      + prometheusQuery.withLegendFormat(|||
        {{%s}} - %s%%
      ||| % [app_label, quantile])
      for quantile in ['50', '90', '99']
    ],

  rateLimitedApps:
    prometheusQuery.new(
      'Prometheus',
      |||
        sum(rate(istio_requests_total{response_flags=~"UO|URX",reporter="source"}[$interval])) by (source_app,destination_service_name)
        /
        sum(rate(istio_requests_total{reporter="source"}[$interval])) by (source_app,destination_service_name)
      |||
    )
    + prometheusQuery.withIntervalFactor(2)
    + prometheusQuery.withLegendFormat('{{source_app}} rate limited to {{destination_service_name}} %'),

  envoyConnections:
    prometheusQuery.new(
      'Prometheus',
      |||
        avg(envoy_server_total_connections{app=~"${sourceApp}"})
      |||
    )
    + prometheusQuery.withIntervalFactor(2)
    + prometheusQuery.withLegendFormat('connections'),


  outgoingRate: rate('source_app'),
  outgoingSuccessRate: succsessRate('source_app'),
  outgoingDuration: duration('source_app'),

  incomingSuccessRate: succsessRate('destination_app'),
  incomingRate: rate('destination_app'),
  incomingDuration: duration('destination_app'),


}

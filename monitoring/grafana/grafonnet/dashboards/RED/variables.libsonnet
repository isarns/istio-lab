local grafonnet = import 'github.com/grafana/grafonnet/gen/grafonnet-latest/main.libsonnet';
local var = grafonnet.dashboard.variable;

{
  new(currentNamespace): {
    datasource:
      var.datasource.new('datasource', 'prometheus')
      + var.datasource.selectionOptions.withMulti(false)
      + var.datasource.generalOptions.withCurrent('Prometheus'),

    namespace:
      var.query.new('namespace')
      + var.query.withDatasourceFromVariable(self.datasource)
      + var.query.queryTypes.withLabelValues(
        'namespace',
        'kube_namespace_status_phase',
      )
      + var.query.selectionOptions.withMulti(false)
      + var.query.generalOptions.withCurrent(currentNamespace)
      + var.query.refresh.onLoad()
      + var.query.refresh.onTime(),

    sourceApp:
      var.query.new('sourceApp')
      + var.query.withDatasourceFromVariable(self.datasource)
      + var.query.queryTypes.withLabelValues(
        'app',
        'istio_requests_total{namespace="$%s"}' % self.namespace.name,
      )
      + var.query.refresh.onLoad()
      + var.query.refresh.onTime(),

    interval:
      var.interval.new('interval', ['1m', '2m', '5m', '10m', '30m', '1h', '6h', '12h', '1d', '7d', '14d', '30d'])
      + var.interval.withAutoOption('30', '2m'),

    deployments_annotation: {
      name: 'Deployment',
      datasource: 'Prometheus',
      showLine: true,
      enable: true,
      expr: 'changes(max(kube_replicaset_created{namespace="$namespace",replicaset=~"$service-.*"})[$__interval:])',
      step: '1m',
      tagKeys: 'deployment',
      titleFormat: 'Deployment',
      iconColor: 'green',
      matchExact: true,
    },

    toArray: [
      self.datasource,
      self.namespace,
      self.sourceApp,
      self.interval,
    ],
  },
}

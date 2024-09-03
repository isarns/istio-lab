local grafonnet = import 'github.com/grafana/grafonnet/gen/grafonnet-latest/main.libsonnet';
local row = grafonnet.panel.row;

local panels = import 'panels.libsonnet';
local variables = import 'variables.libsonnet';
local queries = import 'queries.libsonnet';
local annotations = import 'annotations.libsonnet';

{
  new(): {
    dashboard:
      grafonnet.dashboard.new('RED-Istio')
      + grafonnet.dashboard.withUid('red-istio')
      + grafonnet.dashboard.withDescription('RED for istio')
      + grafonnet.dashboard.withTimezone(value='browser')
      + grafonnet.dashboard.graphTooltip.withSharedCrosshair()
      + grafonnet.dashboard.withTags('red-grafonnet')
      + grafonnet.dashboard.withAnnotations(annotations.toArray)
      + grafonnet.dashboard.withVariables(
        self.variables.toArray
      )
      + grafonnet.dashboard.withPanels(
        grafonnet.util.grid.makeGrid([
          row.new('common')
          + row.withPanels([
            panels.timeSeries.rateLimitedApps('Rate Limited Apps', [queries.rateLimitedApps]),
            panels.timeSeries.envoyConnections('Envoy Connections', [queries.envoyConnections]),
          ]),

          row.new('Incoming ${interval}')
          + row.withPanels([
            panels.timeSeries.requestsPerSecond('Rate', [queries.incomingRate]),
            panels.timeSeries.successRate('Success Rate', [queries.incomingSuccessRate]),
            panels.timeSeries.duration('Duration', queries.incomingDuration),
          ]),
          row.new('Outgoing ${interval}')
          + row.withPanels([
            panels.timeSeries.requestsPerSecond('Rate', [queries.outgoingRate]),
            panels.timeSeries.successRate('Success Rate', [queries.outgoingSuccessRate]),
            panels.timeSeries.duration('Duration', queries.outgoingDuration),
          ]),
        ])
      ),
  },

  withVariables(currentNamespace='apps'): {
    variables:: variables.new(currentNamespace),
  },
}

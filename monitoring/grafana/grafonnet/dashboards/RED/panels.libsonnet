local grafonnet = import 'github.com/grafana/grafonnet/gen/grafonnet-latest/main.libsonnet';

{
  timeSeries: {
    local timeSeries = grafonnet.panel.timeSeries,
    local fieldOverride = grafonnet.panel.timeSeries.fieldOverride,
    local custom = timeSeries.fieldConfig.defaults.custom,
    local options = timeSeries.options,

    base(title, targets):
      timeSeries.new(title)
      + timeSeries.queryOptions.withTargets(targets)
      + options.tooltip.withMode('multi')
      + options.tooltip.withSort('desc')
      + options.legend.withDisplayMode('table')
      + options.legend.withCalcs([
        'lastNotNull',
        'max',
      ])
      + custom.withFillOpacity(10)
      + custom.withShowPoints('never'),

    requestsPerSecond(title, targets):
      timeSeries.new(title)
      + timeSeries.queryOptions.withTargets(targets)
      + timeSeries.standardOptions.withMin(0)
      + timeSeries.standardOptions.withUnit('reqps')
      + timeSeries.standardOptions.withOverrides([
        fieldOverride.byRegexp.new('/active.*/')
        + fieldOverride.byName.withProperty(
          'custom.axisPlacement',
          'right'
        )
        + fieldOverride.byName.withProperty(
          'unit',
          'short'
        ),
      ])
      + options.tooltip.withMode('multi')
      + options.tooltip.withSort('desc')
      + custom.withFillOpacity(10)
      + custom.withLineWidth(2)
      + custom.stacking.withMode('normal'),

    successRate(title, targets):
      timeSeries.new(title)
      + timeSeries.standardOptions.withUnit('percentunit')
      + timeSeries.standardOptions.withDecimals(3)
      + timeSeries.queryOptions.withTargets(targets)
      + options.tooltip.withMode('multi')
      + options.tooltip.withSort('desc')
      + options.legend.withDisplayMode('table')
      + options.legend.withCalcs([
        'lastNotNull',
        'max',
      ])
      + custom.withFillOpacity(10),

    duration(title, targets):
      self.base(title, targets)
      + timeSeries.standardOptions.withUnit('ms')
      + timeSeries.queryOptions.withMaxDataPoints('100')
      + options.tooltip.withMode('multi')
      + options.tooltip.withSort('desc'),

    rateLimitedApps(title, targets):
      self.base(title, targets)
      + timeSeries.standardOptions.withUnit('percentunit'),

    envoyConnections(title, targets):
      self.base(title, targets),
  },


}

// vim: foldmethod=marker foldmarker=local,;

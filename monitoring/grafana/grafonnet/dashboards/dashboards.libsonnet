local RED = import 'RED/dashboard.libsonnet';

{
  'RED.json': (
    RED.new()
    + RED.withVariables(currentNamespace='apps')
  ).dashboard,
}

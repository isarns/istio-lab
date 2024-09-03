{
  deployments: {
    name: 'Deployment',
    datasource: 'Prometheus',
    showLine: true,
    enable: true,
    expr: 'changes(max(kube_replicaset_created{namespace="$namespace",replicaset=~"$sourceApp.*"})[$__interval:])',
    step: '1m',
    tagKeys: 'deployment',
    titleFormat: 'Deployment',
    iconColor: 'green',
    matchExact: true,
  },

  toArray: [
    self.deployments,
  ],
}

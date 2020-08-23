# kubectl-inspection
health check to a kubernetes cluster as a kubectl plugin for ops.

## What's included
- check has no limit workloads
- check notReady nodes

## How to get it
Use go1.13+ and fetch the repo to your any directory using the git command. For example:

```bash
# git clone https://github.com/guanyuv5/inspection-cli-plugin.git
```

## How to use it
```bash
# make
# mv kubectl-inspection /usr/bin/
# kubectl inspection
 Total node count is [3]
 Total notReady node count is [1]
 notReady nodes: [172.21.0.26]
 Total deployment count is [3]
 Total noLimit deployment count is [2]
 noLimit deployment: [default.nginx kube-system.l7-lb-controller]
 Total sts count is [0]
```
# HYPERION

HYPERION is a simple and easy to deploy docker swarm monitoring/alert solution.

- Hyperion's name means "watcher from above" or "he who goes above" from the greek words hyper and iôn.
- Hyperion is tool which provides a very simple way to setup set of FOSS tools to monitor docker swarm cluster.
- Hyperion is for small tech teams who just want a simple and easy way to setup in house monitoring solution for their docker swarm with just one or two commands.
- Hyperion is not fault tolerant (HA Cluster), its a cheap and easy solution, you can configure it accordingly.
- Hyperion is for small swarms with not more than 20 approx nodes.

## *NOTICE : WIP*

# HYPERION Components

- https://github.com/google/cadvisor
- https://github.com/prymitive/karma
- https://github.com/grafana/grafana
- https://github.com/prometheus/prometheus
- https://github.com/prometheus/node_exporter
- https://github.com/prometheus/alertmanager
- https://github.com/containous/traefik

# HYPERION Inspirations

- https://github.com/stefanprodan/swarmprom
- https://github.com/bvis/docker-prometheus-swarm
- https://github.com/ContainerSolutions/prometheus-swarm-discovery
- https://github.com/netman2k/docker-prometheus-swarm
- https://github.com/vegasbrianc/prometheus

# How to run on local machine

```bash
$ docker swarm init
$ docker network create --attachable -d overlay hyperion-net
$ docker stack deploy -c hyperion.yml hyperion
```

# LICENSE

```
This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any
means.

In jurisdictions that recognize copyright laws, the author or authors
of this software dedicate any and all copyright interest in the
software to the public domain. We make this dedication for the benefit
of the public at large and to the detriment of our heirs and
successors. We intend this dedication to be an overt act of
relinquishment in perpetuity of all present and future rights to this
software under copyright law.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

For more information, please refer to <http://unlicense.org>
```
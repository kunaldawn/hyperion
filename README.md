# HYPERION

HYPERION is a simple and easy to deploy docker swarm monitoring/alert solution.

- Hyperion's name means "watcher from above" or "he who goes above" from the greek words hyper and iôn.
- Hyperion is tool which provides a very simple way to setup set of FOSS tools to monitor docker swarm cluster.
- Hyperion is for small tech teams who just want a simple and easy way to setup in house monitoring solution for their docker swarm with just one or two commands.
- Hyperion is not fault tolerant (HA Cluster), its a cheap and easy solution, you can configure it accordingly.
- Hyperion is for small swarms with not more than 20 approx nodes.

*NOTICE : WIP*

# HYPERION Components



# How to run on local machine

```bash
$ docker swarm init
$ docker network create --attachable -d overlay hyperion-net
$ docker stack deploy -c hyperion.yml hyperion
```
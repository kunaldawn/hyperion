# HYPERION

HYPERION was the Titan god of heavenly light. Hyperion's name means "watcher from above" or "he who goes above" from the greek words hyper and iôn.

HYPERION is a docker swarm monitoring/aleart solution.

NOTICE : WIP

## How to run

```bash
$ docker swarm init
$ docker network create --attachable -d overlay hyperion-net
$ docker stack deploy -c hyperion.yml hyperion
```
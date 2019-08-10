#!/bin/sh

docker push kunaldawn/hyperion-prometheus:v1.0;
docker push kunaldawn/hyperion-node-exporter:v1.0;
docker push kunaldawn/hyperion-docker-exporter:v1.0;
docker push kunaldawn/hyperion-alertmanager:v1.0;
docker push kunaldawn/hyperion-grafana:v1.0;
docker push kunaldawn/hyperion-metric-discovery:v1.0;
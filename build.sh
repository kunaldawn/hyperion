#!/bin/sh

cd prometheus && docker build -t kunaldawn/hyperion-prometheus:v1.0 .
cd .. && cd node-exporter && docker build -t kunaldawn/hyperion-node-exporter:v1.0 .
cd .. && cd docker-exporter && docker build -t kunaldawn/hyperion-docker-exporter:v1.0 .
cd .. && cd alertmanager && docker build -t kunaldawn/hyperion-alertmanager:v1.0 .
cd .. && cd grafana && docker build -t kunaldawn/hyperion-grafana:v1.0 .
cd .. && cd metric-discovery && docker build -t kunaldawn/hyperion-metric-discovery:v1.0 .
cd ..
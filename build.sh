#!/bin/sh

cd prometheus && docker build -t kunaldawn/hyperion-prometheus .
cd .. && cd node-exporter && docker build -t kunaldawn/hyperion-node-exporter .
cd .. && cd docker-exporter && docker build -t kunaldawn/hyperion-docker-exporter .
cd .. && cd alertmanager && docker build -t kunaldawn/hyperion-alertmanager .
cd .. && cd grafana && docker build -t kunaldawn/hyperion-grafana .
cd .. && cd metric-discovery && docker build -t kunaldawn/hyperion-metric-discovery .
cd ..
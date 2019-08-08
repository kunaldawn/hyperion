package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

const (
	enableLabel = "hyperion.enable"
	portLabel   = "hyperion.port"
	pathLabel   = "hyperion.path"
)

type MetricProvider struct {
	Targets []string
	Labels  map[string]string
}

type MetricTarget struct {
	Node    swarm.Node
	Service swarm.Service
	Task    swarm.Task
	IP      net.IP
}

type ConnectedTask struct {
	Task swarm.Task
	IP   net.IP
}

type HyperionSwarmDiscovery struct {
	ServiceName       string
	OutputConfigPath  string
	DiscoveryInterval int
}

func NewHyperionSwarmDiscovery(serviceName string, configOutputPath string, discoveryInterval int) *HyperionSwarmDiscovery {
	return &HyperionSwarmDiscovery{ServiceName: serviceName, OutputConfigPath: configOutputPath, DiscoveryInterval: discoveryInterval}
}

func (hyperion *HyperionSwarmDiscovery) getSwarmServiceByName(cli *client.Client, serviceName string) (swarm.Service, error) {
	var service swarm.Service
	serviceFilters := filters.NewArgs()
	serviceFilters.Add("name", serviceName)
	services, err := cli.ServiceList(context.Background(), types.ServiceListOptions{Filters: serviceFilters})
	if err != nil {
		return service, err
	}

	if len(services) == 0 {
		return service, fmt.Errorf("could not find service %s", serviceName)
	}

	service = services[0]
	return service, nil
}

func (hyperion *HyperionSwarmDiscovery) getSwarmServicesByLabel(cli *client.Client, label string) ([]swarm.Service, error) {
	serviceFilters := filters.NewArgs()
	serviceFilters.Add("label", label)
	return cli.ServiceList(context.Background(), types.ServiceListOptions{Filters: serviceFilters})
}

func (hyperion *HyperionSwarmDiscovery) getSwarmServiceTasks(cli *client.Client, serviceID string) ([]swarm.Task, error) {
	taskFilters := filters.NewArgs()
	taskFilters.Add("desired-state", string(swarm.TaskStateRunning))
	taskFilters.Add("service", serviceID)
	return cli.TaskList(context.Background(), types.TaskListOptions{Filters: taskFilters})
}

func (hyperion *HyperionSwarmDiscovery) getSwarmNodes(cli *client.Client) (map[string]swarm.Node, error) {
	nodeMap := make(map[string]swarm.Node)
	nodeFilters := filters.NewArgs()
	nodes, err := cli.NodeList(context.Background(), types.NodeListOptions{Filters: nodeFilters})
	if err != nil {
		return nil, err
	}

	for _, node := range nodes {
		nodeMap[node.ID] = node
	}

	return nodeMap, nil
}

func (hyperion *HyperionSwarmDiscovery) getSwarmServiceNetworks(service swarm.Service) map[string]bool {
	networkIDs := make(map[string]bool)

	for _, virtualIP := range service.Endpoint.VirtualIPs {
		networkIDs[virtualIP.NetworkID] = true
	}

	return networkIDs
}

func (hyperion *HyperionSwarmDiscovery) getSwarmTaskIPs(task swarm.Task) map[string][]net.IP {
	ipsInNetwork := make(map[string][]net.IP)

	for _, networkAttribute := range task.NetworksAttachments {
		if networkAttribute.Network.Spec.Name == "ingress" || networkAttribute.Network.DriverState.Name != "overlay" {
			continue
		}

		ips := make([]net.IP, 0)
		for _, networkCIDR := range networkAttribute.Addresses {
			ip, _, err := net.ParseCIDR(networkCIDR)
			if err != nil {
				continue
			}
			ips = append(ips, ip)
		}

		if len(ips) > 0 {
			ipsInNetwork[networkAttribute.Network.ID] = ips
		}
	}

	return ipsInNetwork
}

func (hyperion *HyperionSwarmDiscovery) getSwarmTasksConnectedToNetwork(tasks []swarm.Task, networkIDs map[string]bool) []ConnectedTask {
	connectedTasks := make([]ConnectedTask, 0)

	for _, task := range tasks {
		ips := hyperion.getSwarmTaskIPs(task)

		for taskNetworkID, taskIPs := range ips {
			if _, ok := networkIDs[taskNetworkID]; ok {
				connectedTasks = append(connectedTasks, ConnectedTask{
					Task: task,
					IP:   taskIPs[0],
				})
			}
		}
	}

	return connectedTasks
}

func (hyperion *HyperionSwarmDiscovery) getMetricTargets(cli *client.Client) ([]MetricTarget, error) {
	prometheusService, err := hyperion.getSwarmServiceByName(cli, hyperion.ServiceName)
	if err != nil {
		return nil, err
	}

	prometheusNetworkIDs := hyperion.getSwarmServiceNetworks(prometheusService)
	prometheusEnabledServices, err := hyperion.getSwarmServicesByLabel(cli, string(enableLabel)+"=true")
	if err != nil {
		return nil, err
	}

	scrapeTargets := make([]MetricTarget, 0)
	nodes, err := hyperion.getSwarmNodes(cli)
	if err != nil {
		return nil, err
	}

	for _, service := range prometheusEnabledServices {
		tasks, err := hyperion.getSwarmServiceTasks(cli, service.ID)
		if err != nil {
			continue
		}

		connectedTasks := hyperion.getSwarmTasksConnectedToNetwork(tasks, prometheusNetworkIDs)
		for _, connectedTask := range connectedTasks {
			target := MetricTarget{
				Node:    nodes[connectedTask.Task.NodeID],
				Service: service,
				Task:    connectedTask.Task,
				IP:      connectedTask.IP,
			}
			scrapeTargets = append(scrapeTargets, target)
		}
	}

	return scrapeTargets, nil
}

func (hyperion *HyperionSwarmDiscovery) buildMetricTargets(scrapeTargets []MetricTarget) []MetricProvider {
	tasks := make([]MetricProvider, 0)
	for _, target := range scrapeTargets {
		task := MetricProvider{
			Targets: hyperion.buildTargets(target),
			Labels:  hyperion.buildLabels(target),
		}
		tasks = append(tasks, task)
	}
	return tasks
}

func (hyperion *HyperionSwarmDiscovery) buildLabels(target MetricTarget) map[string]string {
	labels := map[string]string{
		"job":                        target.Service.Spec.Name,
		"__meta_swarm_task_name":     fmt.Sprintf("%s.%d", target.Service.Spec.Name, target.Task.Slot),
		"__meta_swarm_service_name":  target.Service.Spec.Name,
		"__meta_swarm_node_hostname": target.Node.Description.Hostname,
	}

	if path, ok := target.Service.Spec.Labels[pathLabel]; ok {
		labels["__metrics_path__"] = path
	}

	return labels
}

func (hyperion *HyperionSwarmDiscovery) buildTargets(target MetricTarget) []string {
	var endpoint = target.IP.String()

	if port, ok := target.Service.Spec.Labels[portLabel]; ok {
		endpoint = endpoint + ":" + port
	}

	return []string{endpoint}
}

func (hyperion *HyperionSwarmDiscovery) buildMetricProviderList(cli *client.Client) ([]MetricProvider, error) {
	scrapeTargetsMap, err := hyperion.getMetricTargets(cli)
	if err != nil {
		return nil, err
	}

	return hyperion.buildMetricTargets(scrapeTargetsMap), nil
}

func (hyperion *HyperionSwarmDiscovery) writeMetricProviderConfig(metricProviders []MetricProvider) error {
	jsonData, err := json.MarshalIndent(metricProviders, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(hyperion.OutputConfigPath, jsonData, 0644)
}

func (hyperion *HyperionSwarmDiscovery) sync(cli *client.Client) {
	if _, err := cli.Ping(context.Background()); err != nil {
		log.Println(err)
		return
	}

	providers, err := hyperion.buildMetricProviderList(cli)
	if err != nil {
		log.Println(err)
		return
	}

	if err := hyperion.writeMetricProviderConfig(providers); err != nil {
		log.Println(err)
	} else {
		log.Println(fmt.Sprintf("sync : %d provideros", len(providers)))
	}
}

func (hyperion *HyperionSwarmDiscovery) Run() {
	cli, err := client.NewEnvClient()

	if err != nil {
		log.Panicln(err)
	}

	for {
		hyperion.sync(cli)
		time.Sleep(time.Duration(hyperion.DiscoveryInterval) * time.Second)
	}
}

func main() {
	serviceName, ok := os.LookupEnv("HYPERION_DISCOVERY_SERVICE_NAME")
	if !ok {
		log.Panicln("HYPERION_DISCOVERY_SERVICE_NAME not exported")
	}

	outputConfigPath, ok := os.LookupEnv("HYPERION_DISCOVERY_SERVICE_PATH")
	if !ok {
		log.Panicln("HYPERION_DISCOVERY_SERVICE_PATH not exported")
	}

	refreshIntervalString, ok := os.LookupEnv("HYPERION_DISCOVERY_SERVICE_REFRESH")
	if !ok {
		log.Panicln("HYPERION_DISCOVERY_SERVICE_REFRESH not exported")
	}

	refreshInterval, err := strconv.Atoi(refreshIntervalString)
	if err != nil {
		log.Panicln("HYPERION_DISCOVERY_SERVICE_REFRESH is not int")
	}

	hyperion := NewHyperionSwarmDiscovery(serviceName, outputConfigPath, refreshInterval)
	hyperion.Run()
}

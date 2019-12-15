package main

import (
	"fmt"

	"github.com/VictoriaMetrics/metrics"
)

const (
	namespace = "aws_dc"
)

// NewExporter returns an initialized `Exporter`.
func (hub *Hub) NewExporter(job *Job) (*Exporter, error) {
	dc, err := hub.NewDCClient(&job.AWSCreds)
	if err != nil {
		hub.logger.Errorf("Error initializing AWS Client")
		return nil, err
	}
	return &Exporter{
		client: dc,
		job:    job,
	}, nil
}

// Collect takes in an exporter config, fetches data from AWS APIs
// and constructs metrics for them.
func (hub *Hub) Collect(e Exporter) {
	conn, err := e.client.GetConnections()
	if err != nil {
		hub.logger.Errorf("Error while getting connection states from AWS API: %s", err)
		sendUpMetric(e.job.Name, 0)
		return
	}
	for _, c := range conn.Connections {
		connectionMetricDesc := fmt.Sprintf(`%s_connections{job="%s",conn_state="%s",conn_name="%s",partner_name="%s",conn_id="%s",bandwidth="%s"}`, namespace, e.job.Name, *c.ConnectionState, *c.ConnectionName, *c.PartnerName, *c.ConnectionId, *c.Bandwidth)
		metrics.GetOrCreateGauge(connectionMetricDesc, func() float64 {
			return stateToFloat(*c.ConnectionState)
		})
	}
	interfaces, err := e.client.GetVirtualInterfaces()
	if err != nil {
		hub.logger.Errorf("Error while getting interface states from AWS API: %s", err)
		sendUpMetric(e.job.Name, 0)
		return
	}
	for _, i := range interfaces.VirtualInterfaces {
		interfacesMetricDesc := fmt.Sprintf(`%s_virtual_interfaces{job="%s",virt_interface_state="%s",virt_interface_name="%s",customer_address="%s",virt_interface_id="%s",location="%s"}`, namespace, e.job.Name, *i.VirtualInterfaceState, *i.VirtualInterfaceName, *i.CustomerAddress, *i.VirtualInterfaceId, *i.Location)
		metrics.GetOrCreateGauge(interfacesMetricDesc, func() float64 {
			return stateToFloat(*i.VirtualInterfaceState)
		})

		// fetch list of bgpPeers and create metrics
		for _, bgp := range i.BgpPeers {
			bgpMetricDesc := fmt.Sprintf(`%s_bgp_peers{job="%s",bgp_peer_id="%s",bgp_status="%s",bgp_peer_state="%s",aws_device_v2="%s"}`, namespace, e.job.Name, *bgp.BgpPeerId, *bgp.BgpStatus, *bgp.BgpPeerState, *bgp.AwsDeviceV2)
			metrics.GetOrCreateGauge(bgpMetricDesc, func() float64 {
				return stateToFloat(*bgp.BgpPeerState)
			})
		}
	}
	sendUpMetric(e.job.Name, 1)
}

func sendUpMetric(job string, val float64) {
	upDesc := fmt.Sprintf(`%s_up{job="%s"`, namespace, job)
	metrics.GetOrCreateGauge(upDesc, func() float64 {
		return val
	})
}

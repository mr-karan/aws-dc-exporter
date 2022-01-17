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
func (hub *Hub) Collect(m *metrics.Set, e Exporter) {
	conn, err := e.client.GetConnections()
	upDesc := fmt.Sprintf(`%s_up{job="%s"}`, namespace, e.job.Name)
	if err != nil {
		hub.logger.Errorf("Error while getting connection states from AWS API: %s", err)
		m.GetOrCreateGauge(upDesc, func() float64 {
			return 0
		})
		return
	}
	for _, c := range conn.Connections {
		connectionMetricDesc := fmt.Sprintf(`%s_connections{job="%s",connection_state="%s",connection_name="%s",connection_id="%s",bandwidth="%s",aws_logical_device="%s", location="%s"}`, namespace, e.job.Name, *c.ConnectionState, *c.ConnectionName, *c.ConnectionId, *c.Bandwidth, *c.AwsDeviceV2, *c.Location)
		connState := *c.ConnectionState
		m.GetOrCreateGauge(connectionMetricDesc, func() float64 {
			return stateToFloat(connState)
		})
	}
	interfaces, err := e.client.GetVirtualInterfaces()
	if err != nil {
		hub.logger.Errorf("Error while getting interface states from AWS API: %s", err)
		m.GetOrCreateGauge(upDesc, func() float64 {
			return 0
		})
		return
	}
	for _, i := range interfaces.VirtualInterfaces {
		interfacesMetricDesc := fmt.Sprintf(`%s_virtual_interfaces{job="%s",connection_id="%s",virt_interface_state="%s",virt_interface_name="%s",customer_address="%s",virt_interface_id="%s",aws_logical_device="%s", location="%s"}`, namespace, e.job.Name, *i.ConnectionId, *i.VirtualInterfaceState, *i.VirtualInterfaceName, *i.CustomerAddress, *i.VirtualInterfaceId, *i.AwsDeviceV2, *i.Location)
		intState := *i.VirtualInterfaceState
		m.GetOrCreateGauge(interfacesMetricDesc, func() float64 {
			return stateToFloat(intState)
		})

		// fetch list of bgpPeers and create metrics
		for _, bgp := range i.BgpPeers {
			bgpMetricDesc := fmt.Sprintf(`%s_bgp_peers{job="%s",virt_interface_id="%s",virt_interface_name="%s",bgp_peer_id="%s",bgp_status="%s",bgp_peer_state="%s",aws_logical_device="%s",location="%s"}`, namespace, e.job.Name, *i.VirtualInterfaceId, *i.VirtualInterfaceName, *bgp.BgpPeerId, *bgp.BgpStatus, *bgp.BgpPeerState, *bgp.AwsDeviceV2, *i.Location)
			bgpState := *bgp.BgpPeerState
			m.GetOrCreateGauge(bgpMetricDesc, func() float64 {
				return stateToFloat(bgpState)
			})
		}
	}
	m.GetOrCreateGauge(upDesc, func() float64 {
		return 1
	})
}

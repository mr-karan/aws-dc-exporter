# aws-dc-exporter

## Overview

Export [AWS Direct Connect](https://aws.amazon.com/directconnect/) metrics in Prometheus format.

## Features

- Get metrics of Connections, Virtual Interfaces in form of Prometheus Metrics. (_More metrics coming soon!_)
- Ability to register multiple exporter in form of Jobs to query multiple regions and AWS Accounts.
- Support for `Assume Role` while authenticating to AWS using Role ARN.

## Table of Contents

- [Getting Started](#getting-started)
  - [How it Works](#how-it-works)
  - [Installation](#installation)
  - [Quickstart](#quickstart)
  - [Sending a sample scrape request](#testing-a-sample-alert)

- [Advanced Section](#advanced-section)
  - [Configuration options](#configuation-options)
  - [Setting up Prometheus](#setting-up-prometheus)

## Getting Started

### How it Works

`aws-dc-exporter` uses [AWS SDK](https://github.com/aws/aws-sdk-go) to authenticate with AWS API
and fetch Snapshots metdata. You can specify multiple `jobs` to fetch Direct Connect data and this exporter will collect all metrics and export in the form of Prometheus metrics using a lightweight [metrics](https://github.com/VictoriaMetrics/metrics) collection library.

You will need an _IAM User/Role_ with the following policy attached to the server from where you are running this program:

```plain
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": [
                "directconnect:DescribeConnections",
                "directconnect:DescribeVirtualInterfaces"
            ],
            "Resource": "*"
        }
    ]
}
```

### Installation

There are multiple ways of installing `aws-dc-exporter`.

### Running as docker container

[mrkaran/aws-dc-exporter](https://hub.docker.com/r/mrkaran/aws-dc-exporter)

`docker run -p 9980:9980 -v /etc/aws-dc-exporter/config.toml:/etc/aws-dc-exporter/config.toml mrkaran/aws-dc-exporter:latest`

### Precompiled binaries

Precompiled binaries for released versions are available in the [_Releases_ section](https://github.com/mr-karan/aws-dc-exporter/releases/).

### Compiling the binary

You can checkout the source code and build manually:

```bash
git clone https://github.com/mr-karan/aws-dc-exporter.git
cd aws-dc-exporter
make build
cp config.sample config.toml
./aws-dc-exporter
```

### Quickstart

```sh
mkdir aws-dc-exporter && cd aws-dc-exporter/ # copy the binary and config.sample in this folder
cp config.toml.sample config.toml # change the settings like server address, job metadata, aws credentials etc.
./aws-dc-exporter # this command starts a web server and is ready to collect metrics from EC2.
```

### Testing a sample scrape request

You can send a `GET` request to `/metrics` and see the following metrics in Prometheus format:

```bash
aws_dc_bgp_peers{job="myjob",bgp_peer_id="dxpeer-redacted",bgp_status="up",bgp_peer_state="available",aws_device_v2="xyz-redacted"} 0
aws_dc_connections{job="myjob",conn_state="available",conn_name="redacted",partner_name="xyz",conn_id="dxcon-redacted",bandwidth="100Mbps"} 0
aws_dc_virtual_interfaces{job="myjob",virt_interface_state="available",virt_interface_name="aws-redacted-2",customer_address="x.x.y.z/31",virt_interface_id="dxvif-redacted",location="xyz"} 0
```

## Advanced Section

### Configuration Options

- **[server]**
  - **address**: Port which the server listens to. Default is *9980*
  - **name**: _Optional_, human identifier for the server.
  - **read_timeout**: Duration (in milliseconds) for the request body to be fully read) Read this [blog](https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/) for more info.
  - **write_timeout**: Duration (in milliseconds) for the response body to be written.

- **[app]**
  - **log_level**: "production" for all `INFO` level logs. If you want to enable verbose logging use "debug".
  - **jobs**
    - **name**: Unique identifier for the job.
    - **aws_creds**:
      - **region**: AWS Region where your snapshots are hosted.
      - **access_key**: AWS Access Key if you are using an IAM User. It overrides the env variables `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`.
      - **secret_key**: AWS Secret Key. (See above)
      - **role_arn**: Role ARN if you want to `assume` another role from your IAM Role. This is particularly helpful to scrape data across multiple AWS Accounts.

**NOTE**: You can use `--config` flag to supply a custom config file path while running `aws-dc-exporter`.

### Setting up Prometheus

You can add the following config under `scrape_configs` in Prometheus' configuration.

```yaml
  - job_name: 'aws-dc'
    metrics_path: '/metrics'
    static_configs:
    - targets: ['localhost:9980']
      labels:
        service: direct-connect
```

Validate your setup by querying `aws_dc_up` to check if aws-dc-exporter is discovered by Prometheus:

```plain
`aws_dc_up{job="myjob"} 1`
```

## Example Alerts

<details><summary>Alert when Connection State is not available</summary><br><pre>
- alert: AWSDCConnectionDown
  expr: count(sum(aws_dc_connections{conn_state!="available"}) by (conn_name)) > 0
  for: 1m
  labels:
    room: production-alerts
    severity: warning
  annotations:
    description: AWS Direct Connect Connection {{ $labels.conn_name }} seems to be down.
    title: AWS DC Connection down.
    summary: Please check the AWS DC Console and raise a ticket.
</pre></details>

## Contribution

PRs on Feature Requests, Bug fixes are welcome. Feel free to open an issue and have a discussion first. Contributions on more alert scenarios, more metrics are also welcome and encouraged.

Read [CONTRIBUTING.md](CONTRIBUTING.md) for more details.

## License

[MIT](license)

# Transmission Exporter for Prometheus [![ci](https://github.com/spnngl/transmission-exporter/actions/workflows/ci.yml/badge.svg)](https://github.com/spnngl/transmission-exporter/actions/workflows/ci.yml)

Prometheus exporter for [Transmission](https://transmissionbt.com/) metrics, written in Go.

Fork of [metalmatze/transmission-exporter](https://github.com/metalmatze/transmission-exporter), which is no longer maintained.

### Installation

Each [release](https://github.com/spnngl/transmission-exporter/releases) has a deb and an rpm for amd64 and arm64, with the commands to verify them in its notes. Install one with:

    $ sudo apt install ./transmission-exporter_<version>_amd64.deb
    $ sudo dnf install ./transmission-exporter-<version>-1.x86_64.rpm

Or build it:

    $ go install github.com/spnngl/transmission-exporter/cmd/transmission-exporter@latest

### Configuration

ENV Variable | Description
|----------|-----|
| WEB_PATH | Path for metrics, default: `/metrics` |
| WEB_ADDR | Address for this exporter to run, default: `:19091` |
| TRANSMISSION_ADDR | Transmission address to connect with, default: `http://localhost:9091` |
| TRANSMISSION_USERNAME | Transmission username, no default |
| TRANSMISSION_PASSWORD | Transmission password, no default |

### Container

The image is built in [spnngl/images](https://github.com/spnngl/images).

    docker pull ghcr.io/spnngl/images/transmission-exporter
    docker run -d -p 19091:19091 ghcr.io/spnngl/images/transmission-exporter

### Kubernetes (Prometheus)

A sample kubernetes manifest is available in [examples/kubernetes](examples/kubernetes/transmission.yml)

Please run: `kubectl apply -f examples/kubernetes/transmission.yml`

You should:
* Attach the config and downloads volume
* Configure the password for the exporter

Your prometheus instance will start scraping the metrics automatically. (if configured with annotation based discovery). [more info](https://www.weave.works/docs/cloud/latest/tasks/monitor/configuration-k8s/)

### Docker Compose

[examples/docker-compose.yml](examples/docker-compose.yml) runs Transmission, the exporter and Prometheus:

    docker compose -f examples/docker-compose.yml up

### Development

    make check

Runs every CI check; needs [golangci-lint](https://golangci-lint.run/) on the `PATH` and a C compiler for the race detector. `make install` installs the binary.

Now simply copy the `.env.example` to `.env`, like `cp .env.example .env` and set your preferences.
Now you're good to go.

### Original authors of the Transmission package  
Tobias Blom (https://github.com/tubbebubbe/transmission)  
Long Nguyen (https://github.com/longnguyen11288/go-transmission)

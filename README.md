# `pblb`

[![GoDoc](https://godoc.org/github.com/kineticdial/pblb?status.svg)](https://godoc.org/github.com/kineticdial/pblb)

`pblb` is an experimental load-balancer written in Go. It can't do many things
a production load-balancer can do and—for that reason—should not be used in a
production setting.

## Installation

Clone the repository locally and then `go install` to install the binary on your
path. However, for quick testing, I'd recommend using `docker-compose` (see
"Usage").

## Usage

If you want to quickly test `pblb` I'd recommend using the included
`docker-compose.yaml`. Running `docker-compose up` will set up three NGINX nodes
and the load-balancer configured to these nodes. Once all services are up, you
can hit http://localhost:2839 to see th load-balancer in action, or visit
http://localhost:2840 for prometheus metrics of the load-balancer.

## Configuration

Look at the included `config.yaml` for an example, but I will cover the
highlights here.

- `method`: This field takes a string representation of what load-balancing
method to use. Currently, the only options are `roundrobin` and `twochoice`.
    - [For more information on Round Robin][1].
    - [For more information on Two Choice Random][2].
- `nodes`: This field is an array of dictionaries representing each node to be
load-balanced.
    - `address`: IP/Host Name address to the node.
    - `port`: Port to serve traffic to.
    - `health`: Endpoint to perform health checks on.

[1]: https://www.nginx.com/resources/glossary/round-robin-load-balancing/
[2]: https://www.nginx.com/blog/nginx-power-of-two-choices-load-balancing-algori
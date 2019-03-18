# Vanityd

This is a simple server that allows the `go get` command to retrieve
packages via custom "vanity" URLs. For example, this package can be retrieved
from https://github.com/blinskey/vanityd by running `go get go.linskey.org/vanityd`, and navigating to https://go.linskey.org/vanityd in a browser will bring you to the repository on GitHub via a 302 redirect.

## Basic Installation

- Create an `/etc/vanityd.conf` file defining your settings. A sample file is included in this repository.
- Start the vanityd server.
- Configure your server to forward vanity URL requests to vanityd.

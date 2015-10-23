Installation
============

Build Requirements
------------------

 * Go language environment (>= 1.2)
 * Node.js "npm" package manager
 * "uglifyjs", "uglifycss" and "myth" tools

Debian/Ubuntu:

    apt-get install build-essential golang-go npm nodejs-legacy pandoc
    npm install -g uglifyjs uglifycss myth

Mac OS X (with brew):

    brew install npm pandoc
    npm install -g uglifyjs uglifycss myth

Build Instructions
------------------

Run the building command:

    make
    make install

By default Goverview will be built in the `build` directory and installed in the `/usr/local` one. To change the
installation directory set the `PREFIX` variable:

    sudo make PREFIX=/path/to/directory install

Configuration
-------------

The Goverview configuration consists of a nodes list used by the background poller to query remote monitoring systems.

Main configuration parameters:

 * `bind_addr` _(string)_: service bind address and port
 * `poller` _(object)_: poller configuration

Poller configuration parameters:

 * `nodes` _(list)_: list of nodes to poll

Nodes configuration parameters:

 * `name` _(string)_: unique node name
 * `remote_addr` _(string)_: the remote node address to connect to
 * `label` _(string)_: a user-friendly node name _(optional)_
 * `links` _(list)_: list of access links _(optional)_

Links configuration parameters:

 * `label` _(string)_: link label
 * `host_url` _(string)_: host URL template
 * `service_url` _(string)_: service URL template

In both `host_url` and `service_url`, `%h` and `%s` will be respectively replaced by the host and service names.

Example configuration file:

    ---
    bind_addr: 0.0.0.0:8080
    poller:
      nodes:
        - name: nagios-dc1
          label: Nagios (DC1)
          remote_addr: tcp:192.168.1.10:6557
          poll_interval: 30
          links:
            - label: Go to Nagios
              host_url: 'https://nagios.dc1.example.net/cgi-bin/nagios3/extinfo.cgi?type=1&host=%h'
              service_url: 'https://nagios.dc1.example.net/cgi-bin/nagios3/extinfo.cgi?type=2&host=%h&service=%s'
        - name: centreon-dc2
          label: Host (DC2)
          remote_addr: tcp:192.168.2.10:6557
          poll_interval: 30
           links:
            - label: Go to Centreon
              host_url: 'https://centreon.dc2.example.net/centreon/main.php?p=20102&o=hd&host_name=%h'
              service_url: 'https://centreon.dc2.example.net/centreon/main.php?p=20201&o=svcd&host_name=%h&service_description=%s'


Additional Targets
------------------

Run the various liniting tests:

    make lint

Clean the building environment:

    make clean

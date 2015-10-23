% GOVERVIEW(1) goverview
% Vincent Batoufflet <vincent@batoufflet.info>
% October 23, 2015

# NAME

Goverview - Monitoring overview service

# SYNOPSYS

goverview [*options*]

# DESCRIPTION

Goverview is a monitoring overview service aggregating data from various Nagios compatible sources using Livestatus mechanisms.

# OPTIONS

-c *file*
:   Specify the application configuration file path (type: string, default: /etc/goverview/goverview.yml).

-h
:   Display application help and exit.

-l *file*
:   Specify the server log file (type: string, default: STDERR)

-L *level*
:   Specify the server logging level (type: string, default: info).

    Supported levels: error, warning, notice, info, debug.

-V
:   Display the application version and exit.

# SIGNALS

**goverview** accepts the following signals:

SIGINT, SIGTERM
:   These signals cause **goverview** to terminate.

<https://github.com/vbatoufflet/goverview>

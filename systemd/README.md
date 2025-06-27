# `systemd`
[![Go Reference](https://pkg.go.dev/badge/go.dsage.org/standard/systemd.svg)][1]

This package contains a minimal implementation of the `sd_notify` protocol used to report status about the application
directly to the systemd service manager to both help with automated recovery in the event of errors and give system
administrators more visibility into the service.

These methods are implemented for all unix platforms only, with no-op shims for Windows specifically, and are all safe
to call even when not running as a systemd service.

For more details about the protocol, see [`sd_notify(3)`][2] and [`sd_watchdog_enabled(3)`][3].

[1]: https://pkg.go.dev/go.dsage.org/standard/systemd
[2]: https://www.freedesktop.org/software/systemd/man/latest/sd_notify.html
[3]: https://www.freedesktop.org/software/systemd/man/latest/sd_watchdog_enabled.html

# `systemd`
This package contains a minimal implementation of the `sd_notify` protocol used to report status about the application
directly to the systemd service manager to both help with automated recovery in the event of errors and give system
administrators more visibility into the service.

## Requirements
The functionality in this package requires both of these conditions be met:

1. `GOOS=linux`
2. `NOTIFY_SOCKET` is set in the environment by systemd (service type is `notify` or `notify-reload` in unit file)

If these conditions are not met, all methods in this package are effectively no-ops with no errors. If `NOTIFY_SOCKET`
is set but does not reference a valid `unixgram` socket, an error will be returned.

## Service Status
The [`Notify`][1] and [`Notifyf`][2] methods can be used to report service status directly (see `sd_notify(3)` for valid
states). There are also a variety of convenience methods that can be used to report common states without manually
specifying them.

### Ready
To report the service as ready, call [`NotifyReady`][3].

### Reloading
To report the service is reloading configuration after receiving a signal (usually `SIGHUP`) from systemd, call
[`NotifyReloading`][4]. After the configuration has been reloaded and the service is ready again, call
[`NotifyReady`][3] to tell systemd the reload job has completed.

### Stopping
To report the service is stopping, call [`NotifyStopping`][5].

## Service Watchdog
systemd services can utilize a watchdog to automatically restart the service if it has stopped responding or is no
longer healthy.

### Requirements
In addition to the overall requirements to send notifications, watchdog methods also require an interval is specified in
the `WATCHDOG_USEC` environment variable by systemd (service unit has `WatchdogSec=<secs>`, see `systemd.service(5)`).

### Interval
If enabled for the service, a watchdog interval is configured in the `WATCHDOG_USEC` environment variable. This can be
detected with the [`WatchdogInterval`][7] method.

### Watchdog
Services can implement a watchdog process on their own, or use the [`Watchdog`][6] method for simplicity. This watchdog
process will report to systemd twice during the configured interval and should be run in its own goroutine.

[1]: https://pkg.go.dev/go.dsage.org/standard/systemd#Notify
[2]: https://pkg.go.dev/go.dsage.org/standard/systemd#Notifyf
[3]: https://pkg.go.dev/go.dsage.org/standard/systemd#NotifyReady
[4]: https://pkg.go.dev/go.dsage.org/standard/systemd#NotifyReloading
[5]: https://pkg.go.dev/go.dsage.org/standard/systemd#NotifyStopping
[6]: https://pkg.go.dev/go.dsage.org/standard/systemd#Watchdog
[7]: https://pkg.go.dev/go.dsage.org/standard/systemd#WatchdogInterval

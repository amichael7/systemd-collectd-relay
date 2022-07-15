# systemd-collectd-relay


Relay status data from systemd to collectd using [Go bindings](https://github.com/coreos/go-systemd) to the [dbus API](https://www.freedesktop.org/wiki/Software/systemd/dbus/).  Doesn't make sense to spend lots of time writing a C plugin, an Exec plugin seemed easier.  Pull requests or comments welcome.

## Example Config

```
  LoadPlugin exec
  # ...
  <Plugin exec>
    Exec "root" "/usr/lib/collectd/exec/systemd-collectd-relay" "ssh" "postgresql" "redis"
  </Plugin>
```

# go-cpu-governor

GTK3-based AppIndicator to select a CPU governor. The application consists of 2 parts:

* A back-end service running with root permissions and exporting a D-Bus `system` interface
  `com.github.kamilsamaj.CpuGovernor`, with a single method `com.github.kamilsamaj.CpuGovernor.SetGovernor`.
* GTK3 User-space application displayed as an GTK3 AppIndicator menu, calling
  the `com.github.kamilsamaj.CpuGovernor.SetGovernor` method.

## Build

```shell
make build
```

## Install

```shell
make install
```

## Check the service status

```shell
sudo systemctl status cpu-governor
```

## Check the service logs

```shell
journalctl -fu cpu-governor.service
```

## Check the desktop application logs

```shell
journalctl -f | grep cpu-indicator.desktop
```

## Send a D-Bus message to the service

The back-end service registers a D-Bus `system` interface `com.github.kamilsamaj.CpuGovernor`, with a single
method `SetGovernor`. This method requires a string argument of a value `powersave`, or `performance`.

If the D-Bus call is accepted, the same argument value is returned.

* Set `performance` CPU Governor:

```shell
gdbus call \
    --system \
    --dest com.github.kamilsamaj.CpuGovernor \
    --object-path /com/github/kamilsamaj/CpuGovernor \
    --method com.github.kamilsamaj.CpuGovernor.SetGovernor \
    performance
```

* Set `powersave` CPU Governor:

```shell
gdbus call \
    --system \
    --dest com.github.kamilsamaj.CpuGovernor \
    --object-path /com/github/kamilsamaj/CpuGovernor \
    --method com.github.kamilsamaj.CpuGovernor.SetGovernor \
    powersave
```

## Get the current CPU cores' frequency

```shell
watch cat /sys/devices/system/cpu/cpu[0-9]*/cpufreq/scaling_cur_freq
```

## Get the current CPU governor

```shell
watch cat /sys/devices/system/cpu/cpu[0-9]*/cpufreq/scaling_governor
```

# Troubleshooting

To list and introspect the available D-Bus interfaces, you can install:

```shell
sudo apt install -y qtchooser
```

And run it with:

```shell
qdbusviewer
```

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
method called `SetMode`. This method requires a string argument of a value:

* `powersave`
* `balancepower`
* `balanceperformance`
* `performance`

If the D-Bus call is accepted, the same argument value is returned.

### Using `gdbus` to set the CPU performance and energy preference

```shell
gdbus call \
    --system \
    --dest com.github.kamilsamaj.CpuGovernor \
    --object-path /com/github/kamilsamaj/CpuGovernor \
    --method com.github.kamilsamaj.CpuGovernor.SetMode \
    <REQUESTED_STATE>
```

For example, set the `powersave` CPU Governor with the `power` CPU Power Preference:

```shell
gdbus call \
    --system \
    --dest com.github.kamilsamaj.CpuGovernor \
    --object-path /com/github/kamilsamaj/CpuGovernor \
    --method com.github.kamilsamaj.CpuGovernor.SetMode \
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

# More info

* https://documentation.suse.com/sles/15-SP1/html/SLES-all/cha-tuning-power.html

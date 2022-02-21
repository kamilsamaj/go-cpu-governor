# go-cpu-governor

GTK3-based AppIndicator to select a CPU governor.

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
watch cat /sys/devices/system/cpu/cpu*/cpufreq/scaling_governor
```

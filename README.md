# go-cpu-governor

GTK3-based AppIndicator to select a CPU governor. The application consists of 2 parts:

* A back-end service running with root permissions and exporting a D-Bus `system` interface
  `com.github.kamilsamaj.CpuGovernor`, with a single method `com.github.kamilsamaj.CpuGovernor.SetGovernor`.
* GTK3 User-space application displayed as an GTK3 AppIndicator menu, calling
  the `com.github.kamilsamaj.CpuGovernor.SetGovernor` method.

## Build

### Prerequisites

* Go
* libappindicator-dev 3.0.1
* gir1.2-appindicator3-0.1
* libgtk-3-dev

Ubuntu 22.04+ replaced the original `libappindicator3-dev` with `libayatana-appindicator3-dev` which has a compatible
interface and the existing `.so` objects are symlinked for backwards compatibility but to do the actual compilation of
Go's GTK3 bindings, you'll need to symlink also the C header files and create configuration for the `pkg-config`.

```shell
sudo ln -s /usr/include/libayatana-appindicator3-0.1 /usr/include/libappindicator3-0.1
cd /usr/include/libappindicator3-0.1
sudo ln -s libayatana-appindicator libappindicator

# create a config for pkg-config
cat << 'EOF' | sudo tee /usr/lib/x86_64-linux-gnu/pkgconfig/appindicator3-0.1.pc
prefix=/usr
exec_prefix=${prefix}
libdir=${prefix}/lib/x86_64-linux-gnu
bindir=${exec_prefix}/bin
includedir=${prefix}/include

Cflags: -I${includedir}/libappindicator3-0.1
Requires: dbusmenu-glib-0.4 gtk+-3.0
Libs: -L${libdir} -lappindicator3

Name: appindicator3-0.1
Description: Application indicators
Version: 12.10.0

EOF
```

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

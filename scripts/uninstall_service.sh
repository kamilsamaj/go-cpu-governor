#!/bin/bash
if [ "$UID" -ne 0 ]; then
    echo "$0 can be run only as root" >&2
    exit 1
fi

# disable and stop the service
systemctl stop cpu-governor.service
systemctl disable cpu-governor.service

# stop the desktop app
killall -w cpu-indicator-gtk3

# remove icons
rm -vf /usr/share/icons/hicolor/22x22/status/indicator-cpufreq*.svg
rm -vf /usr/share/icons/ubuntu-mono-dark/status/22/indicator-cpufreq*.svg
rm -vf /usr/share/icons/ubuntu-mono-light/status/22/indicator-cpufreq*.svg
update-icon-caches /usr/share/icons/{ubuntu-mono-dark,ubuntu-mono-light,hicolor}

# remove all configs
rm -vf /etc/dbus-1/system.d/com.github.kamilsamaj.CpuGovernor.conf \
        /usr/lib/systemd/system/cpu-governor.service \
        /etc/xdg/autostart/cpu-indicator.desktop \
        /usr/share/applications/cpu-indicator.desktop \
        /usr/local/bin/cpu-governor-svc \
        /usr/local/bin/cpu-indicator-gtk3

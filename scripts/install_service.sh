#!/bin/bash

CONF_DIR="configs"

if [ "$UID" -ne 0 ]; then
    echo "$0 can be run only as root" >&2
    exit 1
fi

# copy start up configs
cp -vf "${CONF_DIR}/etc/dbus-1/system.d/com.github.kamilsamaj.CpuGovernor.conf" /etc/dbus-1/system.d/
cp -vf "${CONF_DIR}/usr/lib/systemd/system/cpu-governor.service" /usr/lib/systemd/system/
cp -vf "${CONF_DIR}/etc/xdg/autostart/cpu-indicator.desktop" /etc/xdg/autostart/
cp -vf "${CONF_DIR}/etc/xdg/autostart/cpu-indicator.desktop" /usr/share/applications/

# copy the service and application
cp -vf {cpu-governor-svc,cpu-indicator-gtk3} /usr/local/bin/
chmod +x /usr/local/bin/{cpu-governor-svc,cpu-indicator-gtk3}

# enable the service
systemctl enable cpu-governor.service
systemctl start cpu-governor.service

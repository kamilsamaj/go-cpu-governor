#!/bin/bash

CONF_DIR="configs"

if [ "$UID" -ne 0 ]; then
    echo "$0 can be run only as root" >&2
    exit 1
fi

# copy icons
cp -vf "${CONF_DIR}"/usr/share/icons/hicolor/22x22/status/indicator-cpufreq*.svg /usr/share/icons/hicolor/22x22/status/
cp -vf "${CONF_DIR}"/usr/share/icons/ubuntu-mono-dark/status/22/indicator-cpufreq*.svg /usr/share/icons/ubuntu-mono-dark/status/22/
cp -vf "${CONF_DIR}"/usr/share/icons/ubuntu-mono-light/status/22/indicator-cpufreq*.svg /usr/share/icons/ubuntu-mono-light/status/22/
update-icon-caches /usr/share/icons/{ubuntu-mono-dark,ubuntu-mono-light,hicolor}

# copy start-up configs
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

# launch the app
nohup gtk-launch cpu-indicator >/dev/null 2>&1 &

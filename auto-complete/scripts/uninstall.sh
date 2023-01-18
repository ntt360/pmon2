rm -f /etc/init/pmon2.conf
rm -f /etc/logrotate.d/pmon2
rm -f /etc/bash_completion.d/pmon2.sh
systemctl stop pmon2
systemctl disable pmon2
rm -f /usr/lib/systemd/system/pmon2.service
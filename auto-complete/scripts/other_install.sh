systemctl disable pmon2
rm -f /usr/lib/systemd/system/pmon2.service
mkdir -p /usr/lib/systemd/system
cp /usr/local/pmon2/service/centos7/pmon2.service /usr/lib/systemd/system/pmon2.service
cp /usr/local/pmon2/logrotate/pmon2 /etc/logrotate.d/pmon2
cp /usr/local/pmon2/auto-complete/bash/pmon2.sh /etc/bash_completion.d/pmon2.sh
chmod -R 755 /usr/local/pmon2/bin
systemctl enable pmon2
systemctl start pmon2
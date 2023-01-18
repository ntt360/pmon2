echo "请注意：该版本为centos6特制版"
rm -f /etc/init/pmon2.conf
rm -f /etc/init/pmond.conf
mkdir -p "/etc/init/"
cp /usr/local/pmon2/service/centos6/pmon2.conf /etc/init/pmon2.conf
cp /usr/local/pmon2/service/centos6/pmond.conf /etc/init/pmond.conf
cp /usr/local/pmon2/logrotate/pmon2 /etc/logrotate.d/pmon2
cp /usr/local/pmon2/auto-complete/bash/pmon2.sh /etc/bash_completion.d/pmon2.sh
chmod -R 755 /usr/local/pmon2/bin
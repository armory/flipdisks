#!/bin/bash -xe


IP=${1:-"flipdisk.local"}

exit 0
# upload load files
rsync -avzI controller/etc/flipdisk.service "pi@${IP}:/tmp/"
rsync -avzI controller/build/main "pi@${IP}:/home/pi/Desktop/"

ssh pi@${IP} << EOF
  sudo mv /tmp/flipdisk.service /lib/systemd/system/.
  sudo chmod 755 /lib/systemd/system/flipdisk.service
  sudo systemctl enable flipdisk.service
  sudo systemctl start flipdisk
  sudo journalctl -f -u flipdisk
EOF

#!/bin/bash -xe


HOST=${1:-"flipdisk.local"}
PORT=${1:-"22"}

# upload load files
rsync -avzI --port ${PORT} controller/etc/flipdisk.service "pi@${HOST}:/tmp/"
rsync -avzI --port ${PORT} controller/build/main "pi@${HOST}:/home/pi/Desktop/"

ssh -p ${PORT} pi@${HOST} << EOF
  sudo mv /tmp/flipdisk.service /lib/systemd/system/.
  sudo chmod 755 /lib/systemd/system/flipdisk.service
  sudo systemctl enable flipdisk.service
  sudo systemctl restart flipdisk || sudo systemctl start flipdisk
  sudo journalctl -f -u flipdisk
EOF

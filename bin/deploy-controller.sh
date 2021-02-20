#!/bin/bash -xe


HOST=${HOST:-"flipdisk.local"}
PORT=${PORT:-"22"}

# upload load files
rsync -azIv -e "ssh -p${PORT}" controller/etc/flipdisk.service "pi@${HOST}:/tmp/"
rsync -azIv -e "ssh -p${PORT}" controller/build/main "pi@${HOST}:/home/pi/Desktop/"

ssh -p"${PORT}" "pi@${HOST}" << EOF
  sudo mv /tmp/flipdisk.service /lib/systemd/system/.
  sudo chmod 755 /lib/systemd/system/flipdisk.service
  sudo systemctl enable flipdisk.service
  sudo systemctl restart flipdisk || sudo systemctl start flipdisk
  sudo journalctl -f -u flipdisk
EOF

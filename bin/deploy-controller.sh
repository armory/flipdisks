#!/bin/bash -xe

IP=${1-"192.168.86.26"}

# upload load files
rsync -avzI controller/etc/init.d/flipdisk-controller "pi@${IP}:/tmp/"
rsync -avzI controller/build/main "pi@${IP}:/home/pi/Desktop/"


ssh pi@${IP} << EOF
  sudo mv /tmp/flipdisk-controller /etc/init.d/ &&
  sudo chmod +x /etc/init.d/flipdisk-controller &&
  sudo update-rc.d flipdisk-controller defaults &&
  sudo service flipdisk-controller restart &&
  echo "Deployed!"
EOF

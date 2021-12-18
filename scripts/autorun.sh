#!/bin/bash
echo "KILL API"
sudo pkill fibergormapitemplate
sudo systemctl stop fibergormapitemplate

echo "ADD EXECUTION PERMISSIONS"
chmod +x /home/ubuntu/go/src/github.com/nikola43/fibergormapitemplate/fibergormapitemplate

echo "RUN API"
source /home/ubuntu/go/src/github.com/nikola43/fibergormapitemplate/.env
sudo systemctl start fibergormapitemplate
exit

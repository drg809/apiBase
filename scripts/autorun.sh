#!/bin/bash
echo "KILL API"
sudo pkill fibergormapitemplate
sudo systemctl stop fibergormapitemplate

echo "ADD EXECUTION PERMISSIONS"
chmod +x /home/ubuntu/go/src/github.com/drg809/apiBase/fibergormapitemplate

echo "RUN API"
source /home/ubuntu/go/src/github.com/drg809/apiBase/.env
sudo systemctl start fibergormapitemplate
exit

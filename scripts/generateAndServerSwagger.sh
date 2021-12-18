#!/bin/bash
cd ..
rm -rf ./swaggerGenerated/swagger.json
echo "GENERATING..."
swagger generate spec -o ./swaggerGenerated/swagger.json --scan-models

echo "SERVING..."
swagger serve -F=swagger ./swaggerGenerated/swagger.json

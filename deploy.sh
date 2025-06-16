#!/bin/bash

APP_NAME=mlvt
SERVICE_NAME=mlvt
CMD_DIR=cmd/server

echo "🔄 Pulling latest code..."
git pull origin dev 

echo "🛠️ Building binary..."
make build

if [ ! -f $CMD_DIR/$APP_NAME ]; then
    echo "❌ Build failed: Binary not found!"
    exit 1
fi

echo "🚀 Restarting service..."
sudo systemctl restart $SERVICE_NAME

echo "✅ Deployed successfully!"

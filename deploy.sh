#!/bin/bash

export PATH=$PATH:/usr/local/go/bin
echo "PATH: $PATH"
which go
go version

set -e

APP_NAME=mlvt
SERVICE_NAME=mlvt
CMD_DIR=cmd/server
LOG_FILE=/root/code/mlvt-backend/log-deploys/deploy.log

mkdir -p "$(dirname "$LOG_FILE")"

{
  echo ""
  echo "======================"
  echo "🚀 Deploy started at $(date)"
  echo "======================"
  
  cd /root/code/mlvt-backend || exit 1

  echo "🔄 Pulling latest code..."
  git pull origin release/dev 
  echo "✅ Done pulling code"

  echo "🔄 Migrating db up..."
  echo "🔧 Running migrate-up..."
  if ! make migrate-up; then
    echo "❌ migrate-up failed!"
    exit 2
  fi
  echo "✅ Done migrate db up"

  echo "🛠️ Building binary..."
  echo "🛠️ Building binary..."
  if ! make build; then
    echo "❌ Build failed!"
    exit 2
  fi

  if [ ! -f "$CMD_DIR/$APP_NAME" ]; then
      echo "❌ Build failed: Binary not found!"
      exit 1
  fi

  echo "🚀 Restarting service..."
  systemctl restart "$SERVICE_NAME"

  echo "✅ Deployed successfully!"
} >> "$LOG_FILE" 2>&1

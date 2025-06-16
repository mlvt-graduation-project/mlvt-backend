# 🚀 Git Auto Deploy Setup for `mlvt-backend`

This repository is configured for **automatic deployment** using [git-auto-deploy](https://github.com/olipo186/Git-Auto-Deploy). Whenever a webhook is triggered (e.g., from GitHub), the server will automatically pull the latest code, rebuild the binary, and restart the service.

---

## 📦 Requirements

- Ubuntu Server (tested on 20.04+)
- Python 3.10+
- `git-auto-deploy` installed via `pip3 install .` inside the project folder
- A systemd service for your backend (e.g., `mlvt.service`)

---

## ⚙️ Configuration

### 1. **Systemd service: `/etc/systemd/system/git-auto-deploy.service`**

```ini
[Unit]
Description=Git Auto Deploy Webhook Listener
After=network.target

[Service]
ExecStart=/usr/bin/python3 -m gitautodeploy --config /etc/git-auto-deploy.conf --daemon-mode --allow-root-user
Restart=always

[Install]
WantedBy=multi-user.target
```

### 2. **Deployment config: `/etc/git-auto-deploy.conf`**

```json
{
  "repositories": [
    {
      "url": "https://github.com/mlvt-graduation-project/mlvt-backend.git",
      "path": "/root/code/mlvt-backend",
      "deploy": "/root/code/mlvt-backend/deploy.sh"
    }
  ]
}
```

### 3. **Deploy script: `/root/code/mlvt-backend/deploy.sh`**

```bash
#!/bin/bash

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

  echo "🛠️ Building binary..."
  make build

  if [ ! -f "$CMD_DIR/$APP_NAME" ]; then
      echo "❌ Build failed: Binary not found!"
      exit 1
  fi

  echo "🚀 Restarting service..."
  systemctl restart "$SERVICE_NAME"

  echo "✅ Deployed successfully!"
} >> "$LOG_FILE" 2>&1
```

---

## 🚨 Notes

- **Webhook listener** runs on `http://localhost:8001` by default.
- `git-auto-deploy` requires `--allow-root-user` if running under root.
- Make sure `deploy.sh` is executable:
  ```bash
  chmod +x /root/code/mlvt-backend/deploy.sh
  ```

---

## ✅ Testing Deployment

```bash
curl -X POST http://localhost:8001
tail -n 50 /root/code/mlvt-backend/log-deploys/deploy.log
```

---

## 📬 GitHub Webhook

1. Go to your repo → Settings → Webhooks
2. Add webhook:
   - URL: `http://<your-server-ip>:8001`
   - Content type: `application/json`
   - Secret: *(optional)*
   - Trigger: Just `push` events is enough

---

## 💡 Troubleshooting

- If port `8001` is not responding:
  - Check `systemctl status git-auto-deploy`
  - Inspect logs: `journalctl -u git-auto-deploy`
- Make sure your target branch exists in the remote
- Ensure firewall rules allow port `8001` if using public webhooks

---

## ✨ Maintained by Capi

Capybaras don’t rush. Neither does clean code. That’s why Capi’s backend never panics 🛠️
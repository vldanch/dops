# dops — DevOps CLI Assistant

`dops` - is a smart CLI tool for Devops Engineers: checks for the availability of services, log analysis, notifications, system metrics, etc.


## 🚀 Installation

```bash
make build
```

## 🧪 Example usage

Run the ping command to check the availability of one or more URLs concurrently with retries and timeout options:

```bash
./dops ping --url https://example.com --url https://github.com --timeout 5 --retries 3

# Output:
# Ping attempt 1 to https://github.com ... Success! Status code: 200, Response time: 228.728898ms
# Ping attempt 1 to https://example.com ... Success! Status code: 200, Response time: 585.865654ms
```
Run the `checksystem` command to display system metrics:
```bash
./dops checksystem

# Output:
# CPU Usage: 0.00%
# Memory Usage: 3.81% (Used: 298 MB, Total: 7829 MB)
# Disk Usage (/): 0.65% (Used: 6 GB, Total: 1006 GB)
# System Uptime: 3 hours
```

# 📅 ROADMAP

For details on planned features and future development, see the [ROADMAP.md](ROADMAP.md) file.


# 🛠 TODO

## MVP 0.1.0
- ✅ Base project structure
- ✅ `ping` command
- ✅ Reading YAML config

## 0.2.0
- ✅ `checksystem` command
- ❌ Log analyzer

## 0.3.0
- ❌ Telegram notifications
- ❌ `update` command

## Future
- ❌ Web interface
- ❌ SaaS mode



## 💸 Monetization & Licensing

The core of `dops` is and will remain free under the MIT license.

We are planning to introduce a **Pro version** with advanced features, including:

- 🔔 Telegram & Slack notifications
- 📊 Extended system metrics (CPU trends, thresholds)
- 🔄 Auto-update & version check
- 🔐 Role-based config access
- 🌐 Web UI for centralized management (future)

The **Pro version** will be available as a paid binary or Docker container.

## Support

☕ Support my work on coffee: [https://donate.stream/donate_68625552be6ba](https://donate.stream/donate_68625552be6ba)

## Contact

For questions or support, you can reach me on Telegram: [https://t.me/vldanch](https://t.me/vldanch)

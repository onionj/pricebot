# Telegram Price Bot 📊

A Telegram bot that provides real-time updates for currency exchange rates, gold prices, and cryptocurrency values in Iran's market. The bot updates its messages periodically to ensure users always have the latest market data.

## Features 🌟

- Real-time currency exchange rates (USD, EUR, GBP, etc.)
- Gold and coin prices (including Bahar Azadi coin)
- Cryptocurrency prices (Bitcoin, Ethereum, Tether)
- Auto-updating messages
- Persian (Jalali) date support
- State persistence between restarts

## Prerequisites 📋

- Go 1.23.1 or higher
- A Telegram Bot Token
- A Telegram Channel or Chat ID

## Installation 🚀

1. Clone the repository:
   ```bash
   git clone https://github.com/onionj/pricebot.git
   cd pricebot
   ```

2. Set up environment variables:
   ```bash
   cp .env.example .env
   ```
   Edit `.env` and add your:
   - `BOT_TOKEN`: Your Telegram bot token
   - `CHAT_ID`: Target chat/channel ID
   - `CHANEL_NAME`: Your channel name
   - `PROXY_LINK`: (Optional) Proxy link for users

3. Install dependencies:
   ```bash
   go mod download
   ```

## Usage 💡

### Running the Bot

```bash
go run .
```

### Building for Different Platforms

Use the Makefile targets:
```bash
# Build for all platforms
make buildall

# Build for specific platform
make build-linux-amd64
make build-darwin-arm64
```

## Testing 🧪

Run the test suite:
```bash
go test -v ./...
```

Run tests with race condition detection:
```bash
go test -race -v ./...
```

## Quick Install (Linux) 🐧

Install the latest release as a systemd service:
```bash
curl -sfL https://raw.githubusercontent.com/onionj/pricebot/master/install.sh | sudo bash
```

Or from a local binary:
```bash
sudo bash install.sh ./price-linux-amd64
```

This will:
- Install the binary to `/opt/pricebot/`
- Create a `.env` config file (`/opt/pricebot/.env`)
- Set up and enable a `pricebot` systemd service

After installation, edit the config and start the service:
```bash
sudo nano /opt/pricebot/.env
sudo systemctl start pricebot
```

Useful commands:
```bash
sudo systemctl status pricebot      # Check status
sudo journalctl -u pricebot -f      # View logs
sudo systemctl restart pricebot     # Restart
```

## CI/CD 🔄

This project uses GitHub Actions for automated builds and releases. When a version tag is pushed (`v*`):
1. Runs the test suite
2. Builds binaries for all supported platforms (linux/darwin × amd64/arm64)
3. Creates a GitHub Release with binaries and SHA256 checksums

To create a release:
```bash
git tag v1.0.2
git push origin v1.0.2
```

## Project Structure 📁

```
├── .github/workflows/ # GitHub Actions CI/CD
├── price/             # Price fetching and formatting (with retry across multiple endpoints)
├── telegram/          # Telegram bot implementation
├── utils/             # Utility functions (date conversion, etc.)
├── .env.example       # Environment variables template
├── install.sh         # Automated install script (systemd + env)
├── main.go            # Application entry point
└── Makefile           # Build and development commands
```

## Contributing 🤝

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License 📄

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.


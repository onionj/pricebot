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

## CI/CD 🔄

This project uses GitHub Actions for continuous integration. On every push and pull request to the master branch, it:
- Runs the test suite
- Checks for race conditions
- Verifies code builds successfully

## Project Structure 📁

```
├── price/          # Price fetching and formatting
├── telegram/       # Telegram bot implementation
├── utils/          # Utility functions (date conversion, etc.)
├── .env.example    # Environment variables template
├── main.go         # Application entry point
└── Makefile        # Build and development commands
```

## Contributing 🤝

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License 📄

This project is licensed under the MIT License.


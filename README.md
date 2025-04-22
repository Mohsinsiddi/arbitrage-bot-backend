# Arbitrage Bot

A high-performance Ethereum arbitrage bot built in Go.

## Features

- Real-time monitoring of DEX prices
- Automatic detection of arbitrage opportunities
- Optimized transaction submission
- Mempool monitoring
- Gas price optimization
- Multi-node connectivity for reliability

## Getting Started

### Prerequisites

- Go 1.18 or higher
- Ethereum node access (Infura, Alchemy, or local node)
- (Optional) Sepolia testnet ETH for testing

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/arbitrage-bot.git
cd arbitrage-bot

# Install dependencies
go mod download

# Build the application
go build -o arbitrage-bot ./cmd/bot
```

### Configuration

Edit the `configs/config.yaml` file with your settings:

1. Set your Ethereum RPC and WebSocket endpoints
2. Add your wallet's private key (use env vars in production)
3. Configure the contract and token addresses
4. Set DEX router addresses
5. Adjust gas and profit settings

### Running

```bash
# Run with default config
./arbitrage-bot

# Specify a different config file
./arbitrage-bot -config path/to/config.yaml
```

## Architecture

The bot uses a modular architecture with the following components:

- Blockchain communication layer
- Contract interaction
- DEX price monitoring
- Opportunity detection
- Transaction management
- Gas optimization
- Mempool observation

## Development

```bash
# Run tests
go test ./...

# Run linter
golangci-lint run
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
# arbitrage-bot-backend

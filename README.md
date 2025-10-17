# gateway-service-example
Example app running the gateway service for the decentralized XMTP network

## Setup

1. Copy `.env.example` to `.env` and fill in required values:
   - Get Alchemy URLs for the XMTP Ropsten chain (APP_CHAIN) and the Base Sepolia (SETTLEMENT_CHAIN)
   - Generate a payer private key: `xmtpd-cli keys generat
   e`
   - Set the appropriate XMTP environment for your needs

2. Start Redis:
   ```bash
   docker-compose up -d redis
   ```

3. Run the gateway service:
   - Local development: `./dev/start`
   - Docker: `./dev/up`

## Customizing

The Gateway Service is designed to subsidize publishing to the XMTP network. Before deploying publicly, you will want to customize the authorization logic to ensure that only valid requests from your users are allowed through.

Learn more about customizing your Gateway Service [here](https://docs.xmtp.org/fund-agents-apps/run-gateway).

## Installing the xmtpd-cli (optional)

### MacOS
```
brew tap xmtp/tap
brew install xmtpd-cli
```

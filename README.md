# swapica-svc

## Description

Stateless service for signing tx data. Validate crosschain orders.


## Install

  ```bash
  git clone https://github.com/Swapica/swapica-svc
  cd swapica-svc
  go build main.go
  export KV_VIPER_FILE=./config.yaml
  ./main migrate up
  ./main run service
  ```

## Documentation

We do use openapi:json standard for API. We use swagger for documenting our API.

To open online documentation, go to [swagger editor](http://localhost:8080/swagger-editor/) here is how you can start it
```bash
  cd docs
  npm install
  npm start
```
To build documentation use `npm run build` command,
that will create open-api documentation in `web_deploy` folder.

To generate resources for Go models run `./generate.sh` script in root folder.
use `./generate.sh --help` to see all available options.


## Config Setup

```bash
log:
  level: debug
  disable_sentry: true

listener:
  addr: :8000

data:
  tokens: # list of supported tokens
    - id: "1"
      name: "Tether USDT" # token name
      symbol: "USDT" # token symbol
      icon: "https://etherscan.io/token/images/tether_32.png" #token icon
      type: fungible # token type(fungible supports only)
      chains:
        - chain_id: "goerli" 
          token_type: erc20 # token type(erc20/native)
          contract_address: "" # needed for erc20
        - chain_id: "chapel"
          token_type: erc20
          contract_address: ""
  chains:
    - id: "chapel" # chain id for internal search
      name: "BSC Testnet" # chain name
      icon: "https://chainlist.org/_next/image?url=https%3A%2F%2Fdefillama.com%2Fchain-icons%2Frsz_ethereum.jpg&w=64&q=75"
      type: "evm" # chain type(evm supports only)
      swap_contract: "" # address of swap contract
      rpc_endpoint: "" # network rpc
      confirmations: 12 # needed blocks after tx will be confirmed
      chain_params:
        chain_id: 97
        native_symbol: "BNB"
        native_decimals: 18 # native token decimals
        rpc: "https://data-seed-prebsc-1-s1.binance.org:8545" # network rpc
        explorer_url: "https://explorer.binance.org/smart-testnet"
        chain_hex_id: "0x5"
    - id: "goerli"
      name: "Ethereum Goerli"
      icon: "https://chainlist.org/_next/image?url=https%3A%2F%2Fdefillama.com%2Fchain-icons%2Frsz_ethereum.jpg&w=64&q=75"
      type: "evm"
      swap_contract: "" # address of swap contract
      rpc_endpoint: "" # network rpc
      confirmations: 12
      chain_params:
        chain_id: 5
        native_symbol: "ETH"
        native_decimals: 18
        rpc: ""
        explorer_url: "https://goerli.etherscan.io/"
        chain_hex_id: "0x5"

signer:
  eth_signer: "" # Private key of signer

cop:
  disabled: true



```
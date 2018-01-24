
var HDWalletProvider = require("truffle-hdwallet-provider");
var mnemonic = "detect bring exhaust labor liar island dolphin junk total put copy modify";

module.exports = {
  networks: {
    development: {
      host: "localhost",
      port: 8545,
      network_id: "*" // Match any network id
    },
    ropsten: {
      provider: new HDWalletProvider(mnemonic, "https://ropsten.infura.io/aW4LIKVrbtDy5cb8HuMk"),
     port: 8545, 
     network_id:"*",
     gas: 1800000 
    }
  }
};
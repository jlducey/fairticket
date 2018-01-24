
var HDWalletProvider = require("truffle-hdwallet-provider");
var mnemonic = ""detect bring exhaust labor liar island dolphin junk total put copy modify";

module.exports = {
  networks: {
    localhost: {
      host: "localhost", 
      port: 8546, // for ropsten expose this one
      //port: 8545, // expose this one for testrpc
      network_id: "*" // Match any network id
    },  
    ropsten: {
      provider: new HDWalletProvider(mnemonic, 'https://ropsten.infura.io/aW4LIKVrbtDy5cb8HuMk'),
      port: 8545,
      network_id: "*",
      gas: 500000
    }
  }
};
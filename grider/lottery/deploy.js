const HDWalletProvider = require('truffle-hdwallet-provider');
const Web3 = require('web3');
const { interface, bytecode} = require('./compile');
const mnemonic = "candy maple cake sugar pudding cream honey rich smooth crumble sweet treat";
const provider = new HDWalletProvider(mnemonic, "http://localhost:7545");
const web3 = new Web3(provider);

const deploy = async () => {
  const accounts = await web3.eth.getAccounts();
  console.log('Attempting to deploy from account',accounts[0]);
  await new web3.eth.Contract(JSON.parse(interface)
  .deploy({data: bytecode })
  .send({ gas: '100000', from: accounts[0] }));
  console.log(interface);
  console.log('contract deployed to:',result.options.address);
};
deploy();

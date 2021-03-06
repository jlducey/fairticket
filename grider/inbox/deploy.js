const HDWalletProvider = require('truffle-hdwallet-provider');
const Web3 = require('web3');
const { interface, bytecode} = require('./compile');

const provider = new HDWalletProvider(
    'mnemonic for accts', 'https://ropsten.infura.io/aW4LIKVrbtDy5cb8HuMk'
);
const web3 = new Web3(provdier);

const deploy = async () => {
  const accounts = await web3.eth.getAccounts();
  console.log('Attempting to deploy from account',accounts[0]);
  await new web3.eth.Contract(JSON.parse(interface)
  .deploy({data: bytecode, arguments: ['Hi there!']})
  .send({ gas: '100000', from: accounts[0] }));
};
deploy();

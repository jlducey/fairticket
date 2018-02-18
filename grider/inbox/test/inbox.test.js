const assert = require('assert');
const ganache = require('ganache-cli');
const Web3 = require('web3'); // constructor so uppercase on Web3

//const web3 = new Web3(ganache.provider()); // make an instance of web3 tied to ganache
const provider = ganache.provider();
const web3 = new Web3(provider);

const { interface, bytecode } = require('../compile');

let accounts; // let async await ... all async handling
let inbox;

beforeEach( async () => {
// get list of all accounts async
  accounts = await web3.eth.getAccounts();
  inbox = await new web3.eth.Contract(JSON.parse(interface))
  .deploy({ data: bytecode, arguments: ['Hi there!']})
  .send({ from: accounts[0], gas: '1000000' })
  inbox.setProvider(provider);
});

describe('Inbox', ()=> {
    it('deploys a contract', ()=> {
      assert.ok(inbox.options.address);
    });
    it('has a default message', async ()=> {
    const message = await inbox.methods.message().call();
    assert.equal(message,'Hi there!');
    });
    it('can change the message', async ()=> {
    await inbox.methods.setMessage('bye').send({ from: accounts[0] });
    const message = await inbox.methods.message().call();
    assert.equal(message, 'bye');
    });
});
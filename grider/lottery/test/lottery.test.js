// npm run test               in test directory to run this

const assert = require('assert');
const ganache = require('ganache-cli');
const Web3 = require('web3'); // constructor so uppercase on Web3

//const web3 = new Web3(ganache.provider()); // make an instance of web3 tied to ganache
const provider = ganache.provider();
const web3 = new Web3(provider);

const { interface, bytecode } = require('../compile');

// END of generic test stuff.. rest is contract file specific

let accounts; // let async await ... all async handling
let lottery;

beforeEach( async () => {
// get list of all accounts async
  accounts = await web3.eth.getAccounts();
  lottery = await new web3.eth.Contract(JSON.parse(interface))
  .deploy({ data: bytecode })
  .send({ from: accounts[0], gas: '1000000' })
  lottery.setProvider(provider);
});

describe('Lottery Contract', ()=> {
    it('deploys a contract', ()=> {
      assert.ok(lottery.options.address);
    });

    it('allows one account to enter', async () => {
      await lottery.methods.enter().send({
        from: accounts[0],
        value: web3.utils.toWei('0.02','ether')
      });

      const players = await lottery.methods.getPlayers().call({
        from: accounts[0]
      });

      assert.equal(accounts[0], players[0]);
    });

    it('allows multiple accounts to enter', async () => {
      await lottery.methods.enter().send({
        from: accounts[0],
        value: web3.utils.toWei('0.02','ether')
      });

      await lottery.methods.enter().send({
        from: accounts[1],
        value: web3.utils.toWei('0.02','ether')
      });

      await lottery.methods.enter().send({
        from: accounts[2],
        value: web3.utils.toWei('0.02','ether')
      });

      const players = await lottery.methods.getPlayers().call({
        from: accounts[0]
      });

      assert.equal(accounts[0], players[0]);
      assert.equal(accounts[1], players[1]);
      assert.equal(accounts[2], players[2]);
      assert.equal(3, players.length);
    });

      it('requires min .02 ether to enter',async () => {
    try {
      await lottery.methods.enter().send({
      from: accounts[0],
      value: 0
      });
      assert(false);
  } catch (err) {
      assert(err);
  }
    });

    it('only mgr can pick winner', async () => {
      try {
        await lottery.methods.pickWinner().send({
          from: accounts[1]
        });
        assert(false);
      } catch (err) {
        assert(err);
      }
    });





});
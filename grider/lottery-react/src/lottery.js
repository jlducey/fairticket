import web3 from './eb3';

const address = 'address of contract';
const abi = [abi code];

// making local contract instance
export default new web3.eth.Contract(abi, address);
import Web3 from 'web3';

const web3 = new Web3(window.web3.currentProvider);
// above replaces .20 web3 with new 1.x web3
// put console.log(web3.version);   in start of render method of app.js to tell its working
export default web3;
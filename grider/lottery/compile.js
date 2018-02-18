const path = require('path')
const fs = require('fs')
const lotteryPath = path.resolve(__dirname, 'contracts', 'lottery.sol');
const source = fs.readFileSync(lotteryPath, 'utf8');
solc = require('solc'); // added because solc was not found 

module.exports = solc.compile(source,1).contracts[':Lottery'];
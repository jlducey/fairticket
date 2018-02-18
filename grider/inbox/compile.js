const path = require('path')
const fs = require('fs')
const inboxPath = path.resolve(__dirname, 'contracts', 'inbox.sol');
const source = fs.readFileSync(inboxPath, 'utf8');
solc = require('solc'); // added because solc was not found 

module.exports = solc.compile(source,1).contracts[':InBox'];
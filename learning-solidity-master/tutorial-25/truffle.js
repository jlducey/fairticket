module.exports = {
  networks: {
    development: {
      host: "localhost",
      port: 7545,
      network_id: "*" // Match any network id
    },
     ropsten:  {
     network_id: 3,
     host: "localhost",
     port:  7545,
     gas:   2900000
}
  },
   rpc: {
 host: 'localhost',
 post:8080
   }
};
require("@nomicfoundation/hardhat-toolbox");

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: "0.8.24",
  networks: {
    hardhat: {
      mining: {
        auto: true,
        interval: 0,
      },
      hostname: "0.0.0.0",
    },
  },
};

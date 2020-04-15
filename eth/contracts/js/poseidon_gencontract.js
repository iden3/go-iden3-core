#!/usr/bin/env node

const fs = require('fs');

const poseidon = require("./node_modules/circomlib/src/poseidon_gencontract");

let abi = poseidon.abi;
let bytecode = poseidon.createCode();

fs.writeFile("poseidon.abi", JSON.stringify(abi), function(err) {
    if(err) {
        return console.log(err);
    }
});

fs.writeFile("poseidon.bin", bytecode.slice(2), function(err) {
    if(err) {
        return console.log(err);
    }
});


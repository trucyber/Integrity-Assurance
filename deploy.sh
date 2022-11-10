#!/bin/bash

echo "Installing RSV chaincode in peer1"
. set-env.sh peer1 7051
set-chain-env.sh -n rsvCC -v 1.0 -p rsv -I false -C scadachannel
chain.sh install -p

echo "Installing RSV chaincode in peer2"
. set-env.sh peer2 8051
set-chain-env.sh -n rsvCC -v 1.0 -p rsv -I false -C scadachannel
chain.sh install -p

echo "Installing RSV chaincode in peer3"
. set-env.sh peer3 9051
set-chain-env.sh -n rsvCC -v 1.0 -p rsv -I false -C scadachannel
chain.sh install -p

# echo "Installing SPAI chaincode in peer4"
# . set-env.sh peer4 10051
# set-chain-env.sh -n rsvCC -v 1.0 -p rsv -I false -C scadachannel
# chain.sh install -p

# echo "Installing SPAI chaincode in peer5"
# . set-env.sh peer5 11051
# set-chain-env.sh -n rsvCC -v 1.0 -p rsv -I false -C scadachannel
# chain.sh install -p

# echo "Installing SPAI chaincode in peer6"
# . set-env.sh peer6 12051
# set-chain-env.sh -n rsvCC -v 1.0 -p rsv -I false -C scadachannel
# chain.sh install -p

# echo "Installing SPAI chaincode in peer7"
# . set-env.sh peer7 13051
# set-chain-env.sh -n rsvCC -v 1.0 -p rsv -I false -C scadachannel
# chain.sh install -p

# echo "Installing SPAI chaincode in peer8"
# . set-env.sh peer8 14051
# set-chain-env.sh -n rsvCC -v 1.0 -p rsv -I false -C scadachannel
# chain.sh install -p

# echo "Installing SPAI chaincode in peer9"
# . set-env.sh peer9 15051
# set-chain-env.sh -n rsvCC -v 1.0 -p rsv -I false -C scadachannel
# chain.sh install -p

# echo "Installing SPAI chaincode in peer10"
# . set-env.sh peer10 16051
# set-chain-env.sh -n rsvCC -v 1.0 -p rsv -I false -C scadachannel
# chain.sh install -p

# echo "Installing SPAI chaincode in peer11"
# . set-env.sh peer11 17051
# set-chain-env.sh -n rsvCC -v 1.0 -p rsv -I false -C scadachannel
# chain.sh install -p

# echo "Installing SPAI chaincode in peer12"
# . set-env.sh peer12 18051
# set-chain-env.sh -n rsvCC -v 1.0 -p rsv -I false -C scadachannel
# chain.sh install -p

# echo "Installing SPAI chaincode in peer13"
# . set-env.sh peer13 19051
# set-chain-env.sh -n rsvCC -v 1.0 -p rsv -I false -C scadachannel
# chain.sh install -p

# echo "Installing SPAI chaincode in peer14"
# . set-env.sh peer14 20051
# set-chain-env.sh -n rsvCC -v 1.0 -p rsv -I false -C scadachannel
# chain.sh install -p

# echo "Installing SPAI chaincode in peer15"
# . set-env.sh peer15 21051
# set-chain-env.sh -n rsvCC -v 1.0 -p rsv -I false -C scadachannel
# chain.sh install -p

# echo "Installing SPAI chaincode in peer16"
# . set-env.sh peer16 22051
# set-chain-env.sh -n spaiCC -v 1.0 -p spai -I false -C scadachannel
# chain.sh install -p

# echo "Installing SPAI chaincode in peer17"
# . set-env.sh peer17 23051
# set-chain-env.sh -n spaiCC -v 1.0 -p spai -I false -C scadachannel
# chain.sh install -p

# echo "Installing SPAI chaincode in peer18"
# . set-env.sh peer18 24051
# set-chain-env.sh -n spaiCC -v 1.0 -p spai -I false -C scadachannel
# chain.sh install -p

# echo "Installing SPAI chaincode in peer19"
# . set-env.sh peer19 25051
# set-chain-env.sh -n spaiCC -v 1.0 -p spai -I false -C scadachannel
# chain.sh install -p

# echo "Installing SPAI chaincode in peer20"
# . set-env.sh peer20 26051
# set-chain-env.sh -n spaiCC -v 1.0 -p spai -I false -C scadachannel
# chain.sh install -p

# echo "Installing SPAI chaincode in peer21"
# . set-env.sh peer21 27051
# set-chain-env.sh -n spaiCC -v 1.0 -p spai -I false -C scadachannel
# chain.sh install -p

# echo "Installing SPAI chaincode in peer22"
# . set-env.sh peer22 28051
# set-chain-env.sh -n spaiCC -v 1.0 -p spai -I false -C scadachannel
# chain.sh install -p

# echo "Installing SPAI chaincode in peer23"
# . set-env.sh peer23 29051
# set-chain-env.sh -n spaiCC -v 1.0 -p spai -I false -C scadachannel
# chain.sh install -p

# echo "Installing SPAI chaincode in peer24"
# . set-env.sh peer24 30051
# set-chain-env.sh -n spaiCC -v 1.0 -p spai -I false -C scadachannel
# chain.sh install -p

# echo "Installing SPAI chaincode in peer25"
# . set-env.sh peer25 31051
# set-chain-env.sh -n spaiCC -v 1.0 -p spai -I false -C scadachannel
# chain.sh install -p

# echo "Installing SPAI chaincode in peer26"
# . set-env.sh peer26 32051
# set-chain-env.sh -n spaiCC -v 1.0 -p spai -I false -C scadachannel
# chain.sh install -p

# echo "Installing SPAI chaincode in peer27"
# . set-env.sh peer27 33051
# set-chain-env.sh -n spaiCC -v 1.0 -p spai -I false -C scadachannel
# chain.sh install -p

# echo "Installing SPAI chaincode in peer8"
# . set-env.sh peer28 34051
# set-chain-env.sh -n spaiCC -v 1.0 -p spai -I false -C scadachannel
# chain.sh install -p

# echo "Installing SPAI chaincode in peer29"
# . set-env.sh peer29 35051
# set-chain-env.sh -n spaiCC -v 1.0 -p spai -I false -C scadachannel
# chain.sh install -p

# echo "Installing SPAI chaincode in peer30"
# . set-env.sh peer30 36051
# set-chain-env.sh -n spaiCC -v 1.0 -p spai -I false -C scadachannel
# chain.sh install -p

echo    "Instantiating..."
. set-env.sh peer1 7051
set-chain-env.sh -n rsvCC -v 1.0 -p rsv -I false -C scadachannel
# set-chain-env.sh -c   '{"Args":["init","ACFT","1000", "A Cloud Fan Token!!!","john"]}'
chain.sh  instantiate

echo "Done."
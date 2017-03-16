#/usr/bin/bash
./simple-raft -id 1 -port :50001
./simple-raft -id 2 -port :50002
./simple-raft -id 3 -port :50003
./simple-raft -id 4 -port :50004
./simple-raft -id 0 -port :50000 -state 0

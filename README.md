# simple blockchain thingy

## testing
go test -v -run TestPrepareData
go test -v -run TestRun

## boltDB
In blocks, the key -> value pairs are:

'b' + 32-byte block hash -> block index record
'f' + 4-byte file number -> file information record
// we build a file number from the block height by dividing by 1000 and adding 1
'l' -> 4-byte file number: the last block file number used
'R' -> 1-byte boolean: whether we're in the process of reindexing
'F' + 1-byte flag name length + flag name string -> 1 byte boolean: various flags that can be on or off
't' + 32-byte transaction hash -> transaction index record
In chainstate, the key -> value pairs are:

'c' + 32-byte transaction hash -> unspent transaction output record for that transaction
'B' -> 32-byte block hash: the block hash up to which the database represents the unspent transaction outputs
(Detailed explanation can be found here)

https://en.bitcoin.it/wiki/Bitcoin_Core_0.11_(ch_2):_Data_Storage
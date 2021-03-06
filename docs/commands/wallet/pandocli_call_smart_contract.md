## pandocli call smart_contract

Call or deploy a smart contract

### Synopsis

smartContractCmd represents the smart_contract command, which can be used to calls the specified smart contract.
		However, calling a smart contract does NOT modify the globally consensus state. It can be used for dry run, or for retrieving info from smart contracts without actually spending gas.

```
pandocli call smart_contract [flags]
```

### Examples

```

	[Deploy a smart contract (local only)]
	pandocli call smart_contract --from=df1f3D3eE9430dB3A44aE6B80Eb3E23352BB785E --value=1680 --gas_price=3 --gas_limit=50000 --data=600a600c600039600a6000f3600360135360016013f3
	
	[Call an API of a smart contract (local only)]
	pandocli call smart_contract --from=df1f3D3eE9430dB3A44aE6B80Eb3E23352BB785E --to=0x7ad6cea2bc3162e30a3c98d84f821b3233c22647 --gas_price=3 --gas_limit=50000
	
```

### Options

```
      --chain string       Chain ID
      --data string        The data for the smart contract
      --from string        The caller address
      --gas_limit uint     The gas limit
      --gas_price string   The gas price (default "100000000wei")
  -h, --help               help for smart_contract
      --seq uint           Sequence number of the transaction
      --to string          The smart contract address
      --value uint         Value to be transferred
```

### Options inherited from parent commands

```
      --config string   config path (default is /Users/<username>/.pandocli) (default "/Users/<username>/.pandocli")
```

### SEE ALSO

* [pandocli call](pandocli_call.md)	 - Call smart contract APIs

###### Auto generated by spf13/cobra on 24-Apr-2019

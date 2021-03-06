package call

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	rpcc "github.com/ybbus/jsonrpc"

	"github.com/pandotoken/pando/cmd/pandocli/cmd/utils"
	"github.com/pandotoken/pando/common"
	"github.com/pandotoken/pando/ledger/types"
	"github.com/pandotoken/pando/rpc"
)

// smartContractCmd represents the smart_contract command, which can be used to calls the specified smart contract.
// However, calling a smart contract does NOT modify the globally consensus state. It can be used for dry run, or
// for retrieving info from smart contracts without actually spending gas.
// Examples:
//   * Deploy a smart contract (local only)
//		pandocli call smart_contract --from=df1f3D3eE9430dB3A44aE6B80Eb3E23352BB785E --value=1680 --gas_price=3 --gas_limit=50000 --data=600a600c600039600a6000f3600360135360016013f3
//   * Call an API of a smart contract (local only)
//		pandocli call smart_contract --from=df1f3D3eE9430dB3A44aE6B80Eb3E23352BB785E --to=0x7ad6cea2bc3162e30a3c98d84f821b3233c22647 --gas_price=3 --gas_limit=50000

var smartContractCmd = &cobra.Command{
	Use:   "smart_contract",
	Short: "Call or deploy a smart contract",
	Example: `
	[Deploy a smart contract (local only)]
	pandocli call smart_contract --from=df1f3D3eE9430dB3A44aE6B80Eb3E23352BB785E --value=1680 --gas_price=3 --gas_limit=50000 --data=600a600c600039600a6000f3600360135360016013f3
	
	[Call an API of a smart contract (local only)]
	pandocli call smart_contract --from=df1f3D3eE9430dB3A44aE6B80Eb3E23352BB785E --to=0x7ad6cea2bc3162e30a3c98d84f821b3233c22647 --gas_price=3 --gas_limit=50000
	`,
	Long: `smartContractCmd represents the smart_contract command, which can be used to calls the specified smart contract.
		However, calling a smart contract does NOT modify the globally consensus state. It can be used for dry run, or for retrieving info from smart contracts without actually spending gas.`,
	Run: doSmartContractCmd,
}

func doSmartContractCmd(cmd *cobra.Command, args []string) {
	from := types.TxInput{
		Address: common.HexToAddress(fromFlag),
		Coins: types.Coins{
			PandoWei: new(big.Int).SetUint64(0),
			PTXWei:   new(big.Int).SetUint64(valueFlag),
		},
		Sequence: seqFlag,
	}

	to := types.TxOutput{
		Address: common.HexToAddress(toFlag),
	}

	gasPrice, ok := types.ParseCoinAmount(gasPriceFlag)
	if !ok {
		utils.Error("Failed to parse gas price")
	}

	data, err := hex.DecodeString(dataFlag)
	if err != nil {
		utils.Error("Failed to decode data: %v, err: %v\n", dataFlag, err)
	}

	sctx := &types.SmartContractTx{
		From:     from,
		To:       to,
		GasLimit: gasLimitFlag,
		GasPrice: gasPrice,
		Data:     data,
	}

	sctxBytes, err := types.TxToBytes(sctx)
	if err != nil {
		utils.Error("Failed to encode smart contract transaction: %v\n", sctx)
	}
	if verboseFlag {
		fmt.Printf("Encoded Tx: %x\n\n", sctxBytes)
	}

	rpcCallArgs := rpc.CallSmartContractArgs{
		SctxBytes: hex.EncodeToString(sctxBytes),
	}

	client := rpcc.NewRPCClient(viper.GetString(utils.CfgRemoteRPCEndpoint))

	res, err := client.Call("pando.CallSmartContract", rpcCallArgs)
	if err != nil {
		utils.Error("Failed to call smart contract: %v\n", err)
	}
	if res.Error != nil {
		utils.Error("Failed to execute smart contract: %v\n", res.Error)
	}
	json, err := json.MarshalIndent(res.Result, "", "    ")
	if err != nil {
		utils.Error("Failed to parse server response: %v\n%s\n", err, string(json))
	}
	fmt.Println(string(json))
}

func init() {
	smartContractCmd.Flags().StringVar(&chainIDFlag, "chain", "", "Chain ID")
	smartContractCmd.Flags().StringVar(&fromFlag, "from", "", "The caller address")
	smartContractCmd.Flags().StringVar(&toFlag, "to", "", "The smart contract address")
	smartContractCmd.Flags().Uint64Var(&valueFlag, "value", 0, "Value to be transferred")
	smartContractCmd.Flags().StringVar(&gasPriceFlag, "gas_price", fmt.Sprintf("%dwei", types.MinimumGasPrice), "The gas price")
	smartContractCmd.Flags().Uint64Var(&gasLimitFlag, "gas_limit", 0, "The gas limit")
	smartContractCmd.Flags().StringVar(&dataFlag, "data", "", "The data for the smart contract")
	smartContractCmd.Flags().Uint64Var(&seqFlag, "seq", 0, "Sequence number of the transaction")
	smartContractCmd.Flags().BoolVar(&verboseFlag, "verbose", false, "")

	smartContractCmd.MarkFlagRequired("from")
	smartContractCmd.MarkFlagRequired("gas_price")
	smartContractCmd.MarkFlagRequired("gas_limit")
}

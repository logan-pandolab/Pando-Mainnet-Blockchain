package tx

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	"github.com/pandotoken/pando/cmd/pandocli/cmd/utils"
	"github.com/pandotoken/pando/common"
	"github.com/pandotoken/pando/ledger/types"
	"github.com/pandotoken/pando/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	rpcc "github.com/ybbus/jsonrpc"
)

// splitRuleCmd represents the split rule command
// Example:
//		pandocli tx split_rule --chain="pandonet" --from=df1f3D3eE9430dB3A44aE6B80Eb3E23352BB785E --seq=8 --resource_id=die_another_day --addresses=df1f3D3eE9430dB3A44aE6B80Eb3E23352BB785E,df1f3D3eE9430dB3A44aE6B80Eb3E23352BB785E --percentages=30,30 --duration=1000
var splitRuleCmd = &cobra.Command{
	Use:     "split_rule",
	Short:   "Initiate or update a split rule",
	Example: `pandocli tx split_rule --chain="pandonet" --from=df1f3D3eE9430dB3A44aE6B80Eb3E23352BB785E --seq=8 --resource_id=die_another_day --addresses=df1f3D3eE9430dB3A44aE6B80Eb3E23352BB785E,df1f3D3eE9430dB3A44aE6B80Eb3E23352BB785E --percentages=30,30 --duration=1000`,
	Run:     doSplitRuleCmd,
}

func doSplitRuleCmd(cmd *cobra.Command, args []string) {
	wallet, fromAddress, err := walletUnlock(cmd, fromFlag)
	if err != nil {
		return
	}
	defer wallet.Lock(fromAddress)

	input := types.TxInput{
		Address:  fromAddress,
		Sequence: uint64(seqFlag),
	}

	if len(addressesFlag) != len(percentagesFlag) {
		fmt.Println("Should have the same number of addresses and percentages")
		return
	}
	var splits []types.Split
	for idx, addressStr := range addressesFlag {
		percentageStr := percentagesFlag[idx]

		address, err := hex.DecodeString(addressStr)
		if err != nil {
			fmt.Println("The address must be a hex string")
			return
		}

		percentage, err := strconv.ParseUint(percentageStr, 10, 32)
		if err != nil {
			fmt.Println(err)
			return
		}

		split := types.Split{
			Address:    common.BytesToAddress(address),
			Percentage: uint(percentage),
		}
		splits = append(splits, split)
	}

	fee, ok := types.ParseCoinAmount(feeFlag)
	if !ok {
		utils.Error("Failed to parse fee")
	}

	splitRuleTx := &types.SplitRuleTx{
		Fee: types.Coins{
			PandoWei: new(big.Int).SetUint64(0),
			PTXWei:   fee,
		},
		ResourceID: resourceIDFlag,
		Initiator:  input,
		Duration:   durationFlag,
		Splits:     splits,
	}

	sig, err := wallet.Sign(fromAddress, splitRuleTx.SignBytes(chainIDFlag))
	if err != nil {
		utils.Error("Failed to sign transaction: %v\n", err)
	}
	splitRuleTx.SetSignature(fromAddress, sig)

	raw, err := types.TxToBytes(splitRuleTx)
	if err != nil {
		utils.Error("Failed to encode transaction: %v\n", err)
	}
	signedTx := hex.EncodeToString(raw)

	client := rpcc.NewRPCClient(viper.GetString(utils.CfgRemoteRPCEndpoint))

	res, err := client.Call("pando.BroadcastRawTransaction", rpc.BroadcastRawTransactionArgs{TxBytes: signedTx})
	if err != nil {
		utils.Error("Failed to broadcast transaction: %v\n", err)
	}
	if res.Error != nil {
		utils.Error("Server returned error: %v\n", res.Error)
	}
	fmt.Printf("Successfully broadcasted transaction.\n")
}

func init() {
	splitRuleCmd.Flags().StringVar(&chainIDFlag, "chain", "", "Chain ID")
	splitRuleCmd.Flags().StringVar(&fromFlag, "from", "", "Initiator's address")
	splitRuleCmd.Flags().Uint64Var(&seqFlag, "seq", 0, "Sequence number of the transaction")
	splitRuleCmd.Flags().StringVar(&feeFlag, "fee", fmt.Sprintf("%dwei", types.MinimumTransactionFeePTXWei), "Fee")
	splitRuleCmd.Flags().StringVar(&resourceIDFlag, "resource_id", "", "The resourceID of interest")
	splitRuleCmd.Flags().StringSliceVar(&addressesFlag, "addresses", []string{}, "List of addresses participating in the split")
	splitRuleCmd.Flags().StringSliceVar(&percentagesFlag, "percentages", []string{}, "List of integers (between 0 and 100) representing of percentage of split")
	splitRuleCmd.Flags().Uint64Var(&durationFlag, "duration", 1000, "Reserve duration")
	splitRuleCmd.Flags().StringVar(&walletFlag, "wallet", "soft", "Wallet type (soft|nano)")

	splitRuleCmd.MarkFlagRequired("chain")
	splitRuleCmd.MarkFlagRequired("from")
	splitRuleCmd.MarkFlagRequired("seq")
	splitRuleCmd.MarkFlagRequired("addresses")
	splitRuleCmd.MarkFlagRequired("percentages")
	splitRuleCmd.MarkFlagRequired("resource_id")
	splitRuleCmd.MarkFlagRequired("duration")
}

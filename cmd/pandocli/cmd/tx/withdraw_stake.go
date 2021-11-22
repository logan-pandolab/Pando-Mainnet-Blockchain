package tx

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/pandotoken/pando/cmd/pandocli/cmd/utils"
	"github.com/pandotoken/pando/common"
	"github.com/pandotoken/pando/ledger/types"
	"github.com/pandotoken/pando/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	rpcc "github.com/ybbus/jsonrpc"
)

// withdrawStakeCmd represents the withdraw stake command
// Example:
//		pandocli tx withdraw --chain="pandonet" --source=df1f3D3eE9430dB3A44aE6B80Eb3E23352BB785E --holder=df1f3D3eE9430dB3A44aE6B80Eb3E23352BB785E --purpose=0 --seq=8
var withdrawStakeCmd = &cobra.Command{
	Use:     "withdraw",
	Short:   "withdraw stake to a validator or guardian",
	Example: `pandocli tx withdraw --chain="pandonet" --source=df1f3D3eE9430dB3A44aE6B80Eb3E23352BB785E --holder=df1f3D3eE9430dB3A44aE6B80Eb3E23352BB785E --purpose=0 --seq=8`,
	Run:     doWithdrawStakeCmd,
}

func doWithdrawStakeCmd(cmd *cobra.Command, args []string) {
	wallet, sourceAddress, err := walletUnlockWithPath(cmd, sourceFlag, pathFlag)
	if err != nil {
		return
	}
	defer wallet.Lock(sourceAddress)

	fee, ok := types.ParseCoinAmount(feeFlag)
	if !ok {
		utils.Error("Failed to parse fee")
	}

	source := types.TxInput{
		Address:  sourceAddress,
		Sequence: uint64(seqFlag),
	}
	holder := types.TxOutput{
		Address: common.HexToAddress(holderFlag),
	}

	withdrawStakeTx := &types.WithdrawStakeTx{
		Fee: types.Coins{
			PandoWei: new(big.Int).SetUint64(0),
			PTXWei:   fee,
		},
		Source:  source,
		Holder:  holder,
		Purpose: purposeFlag,
	}

	sig, err := wallet.Sign(sourceAddress, withdrawStakeTx.SignBytes(chainIDFlag))
	if err != nil {
		utils.Error("Failed to sign transaction: %v\n", err)
	}
	withdrawStakeTx.SetSignature(sourceAddress, sig)

	raw, err := types.TxToBytes(withdrawStakeTx)
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
	withdrawStakeCmd.Flags().StringVar(&chainIDFlag, "chain", "", "Chain ID")
	withdrawStakeCmd.Flags().StringVar(&sourceFlag, "source", "", "Source of the stake")
	withdrawStakeCmd.Flags().StringVar(&holderFlag, "holder", "", "Holder of the stake")
	withdrawStakeCmd.Flags().StringVar(&pathFlag, "path", "", "Wallet derivation path")
	withdrawStakeCmd.Flags().StringVar(&feeFlag, "fee", fmt.Sprintf("%dwei", types.MinimumTransactionFeePTXWei), "Fee")
	withdrawStakeCmd.Flags().Uint64Var(&seqFlag, "seq", 0, "Sequence number of the transaction")
	withdrawStakeCmd.Flags().Uint8Var(&purposeFlag, "purpose", 0, "Purpose of staking")
	withdrawStakeCmd.Flags().StringVar(&walletFlag, "wallet", "soft", "Wallet type (soft|nano)")

	withdrawStakeCmd.MarkFlagRequired("chain")
	withdrawStakeCmd.MarkFlagRequired("source")
	withdrawStakeCmd.MarkFlagRequired("holder")
	withdrawStakeCmd.MarkFlagRequired("seq")
}

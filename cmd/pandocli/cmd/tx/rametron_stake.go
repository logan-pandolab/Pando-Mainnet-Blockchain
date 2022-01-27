package tx

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/pandotoken/pando/cmd/pandocli/cmd/utils"
	"github.com/pandotoken/pando/common"
	"github.com/pandotoken/pando/ledger/types"
	"github.com/pandotoken/pando/rpc"
	wtypes "github.com/pandotoken/pando/wallet/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ybbus/jsonrpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// rametronStakeCmd represents the rametronStake command
// Example:
//		pandocli tx rametronStake --chain="pandonet" --from=df1f3D3eE9430dB3A44aE6B80Eb3E23352BB785E --to=df1f3D3eE9430dB3A44aE6B80Eb3E23352BB785E --pando=10 --ptx=9 --seq=1
//		pandocli tx rametronStake --chain="pandonet" --path "m/44'/60'/0'/0/0" --to=df1f3D3eE9430dB3A44aE6B80Eb3E23352BB785E --pando=10 --ptx=9 --seq=1 --wallet=trezor
//		pandocli tx rametronStake --chain="pandonet" --path "m/44'/60'/0'/0" --to=df1f3D3eE9430dB3A44aE6B80Eb3E23352BB785E --pando=10 --ptx=9 --seq=1 --wallet=nano
var rametronStakeCmd = &cobra.Command{
	Use:     "rametronStake",
	Short:   "RametronStake tokens",
	Example: `pandocli tx rametronStake --chain="pandonet" --from=df1f3D3eE9430dB3A44aE6B80Eb3E23352BB785E --to=df1f3D3eE9430dB3A44aE6B80Eb3E23352BB785E --pando=10 --ptx=9 --seq=1`,
	Run:     doRametronStakeCmd,
}

func doRametronStakeCmd(cmd *cobra.Command, args []string) {
	walletType := getWalletType(cmd)
	if walletType == wtypes.WalletTypeSoft && len(fromFlag) == 0 {
		utils.Error("The from address cannot be empty") // we don't need to specify the "from address" for hardware wallets
		return
	}

	if len(toFlag) == 0 {
		utils.Error("The to address cannot be empty")
		return
	}
	if fromFlag == toFlag {
		utils.Error("The from and to address cannot be identical")
		return
	}

	wallet, fromAddress, err := walletUnlockWithPath(cmd, fromFlag, pathFlag)
	if err != nil || wallet == nil {
		return
	}
	defer wallet.Lock(fromAddress)

	pando, ok := types.ParseCoinAmount(pandoAmountFlag)
	if !ok {
		utils.Error("Failed to parse pando amount")
	}
	ptx, ok := types.ParseCoinAmount(ptxAmountFlag)
	if !ok {
		utils.Error("Failed to parse ptx amount")
	}
	fee, ok := types.ParseCoinAmount(feeFlag)
	if !ok {
		utils.Error("Failed to parse fee")
	}
	inputs := []types.TxInput{{
		Address: fromAddress,
		Coins: types.Coins{
			PTXWei:   new(big.Int).Add(ptx, fee),
			PandoWei: pando,
		},
		Sequence: uint64(seqFlag),
	}}
	outputs := []types.TxOutput{{
		Address: common.HexToAddress(toFlag),
		Coins: types.Coins{
			PTXWei:   ptx,
			PandoWei: pando,
		},
	}}
	rametronStakeTx := &types.RametronStakeTx{
		Fee: types.Coins{
			PandoWei: new(big.Int).SetUint64(0),
			PTXWei:   fee,
		},
		Inputs:  inputs,
		Outputs: outputs,
	}

	sig, err := wallet.Sign(fromAddress, rametronStakeTx.SignBytes(chainIDFlag))
	if err != nil {
		utils.Error("Failed to sign transaction: %v\n", err)
	}
	rametronStakeTx.SetSignature(fromAddress, sig)

	raw, err := types.TxToBytes(rametronStakeTx)
	if err != nil {
		utils.Error("Failed to encode transaction: %v\n", err)
	}
	signedTx := hex.EncodeToString(raw)

	client := rpcc.NewRPCClient(viper.GetString(utils.CfgRemoteRPCEndpoint))

	var res *jsonrpc.RPCResponse
	if asyncFlag {
		res, err = client.Call("pando.BroadcastRawTransactionAsync", rpc.BroadcastRawTransactionArgs{TxBytes: signedTx})
	} else {
		res, err = client.Call("pando.BroadcastRawTransaction", rpc.BroadcastRawTransactionArgs{TxBytes: signedTx})
	}

	if err != nil {
		utils.Error("Failed to broadcast transaction: %v\n", err)
	}
	if res.Error != nil {
		utils.Error("Server returned error: %v\n", res.Error)
	}
	result := &rpc.BroadcastRawTransactionResult{}
	err = res.GetObject(result)
	if err != nil {
		utils.Error("Failed to parse server response: %v\n", err)
	}
	formatted, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		utils.Error("Failed to parse server response: %v\n", err)
	}
	fmt.Printf("Successfully broadcasted transaction:\n%s\n", formatted)
}

func init() {
	rametronStakeCmd.Flags().StringVar(&chainIDFlag, "chain", "", "Chain ID")
	rametronStakeCmd.Flags().StringVar(&fromFlag, "from", "", "Source of the stake")
	rametronStakeCmd.Flags().StringVar(&toFlag, "to", "", "Holder of the stake")
	rametronStakeCmd.Flags().StringVar(&pathFlag, "path", "", "Wallet derivation path")
	rametronStakeCmd.Flags().Uint64Var(&seqFlag, "seq", 0, "Sequence number of the transaction")
	rametronStakeCmd.Flags().StringVar(&pandoAmountFlag, "pando", "0", "Pando amount")
	rametronStakeCmd.Flags().StringVar(&ptxAmountFlag, "ptx", "0", "Pando amount")
	rametronStakeCmd.Flags().StringVar(&feeFlag, "fee", fmt.Sprintf("%dwei", types.MinimumTransactionFeePTXWei), "Fee")
	rametronStakeCmd.Flags().StringVar(&walletFlag, "wallet", "soft", "Wallet type (soft|nano|trezor)")
	rametronStakeCmd.Flags().BoolVar(&asyncFlag, "async", false, "block until tx has been included in the blockchain")

	rametronStakeCmd.MarkFlagRequired("chain")
	//rametronStakeCmd.MarkFlagRequired("from")
	rametronStakeCmd.MarkFlagRequired("to")
	rametronStakeCmd.MarkFlagRequired("seq")
}


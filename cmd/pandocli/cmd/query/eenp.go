package query

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/pandotoken/pando/cmd/pandocli/cmd/utils"
	"github.com/pandotoken/pando/common"
	"github.com/pandotoken/pando/rpc"

	rpcc "github.com/ybbus/jsonrpc"
)

// eenpCmd represents the eenp command.
// Example:
//		pandocli query eenp --height=10
var eenpCmd = &cobra.Command{
	Use:     "eenp",
	Short:   "Get elite edge node pool",
	Example: `pandocli query eenp --height=10`,
	Run:     doEenpCmd,
}

func doEenpCmd(cmd *cobra.Command, args []string) {
	client := rpcc.NewRPCClient(viper.GetString(utils.CfgRemoteRPCEndpoint))

	height := heightFlag
	res, err := client.Call("pando.GetEenpByHeight", rpc.GetEenpByHeightArgs{Height: common.JSONUint64(height)})
	if err != nil {
		utils.Error("Failed to get elite edge node pool: %v\n", err)
	}
	if res.Error != nil {
		utils.Error("Failed to get elite edge node pool: %v\n", res.Error)
	}
	json, err := json.MarshalIndent(res.Result, "", "    ")
	if err != nil {
		utils.Error("Failed to parse server response: %v\n%s\n", err, string(json))
	}
	fmt.Println(string(json))
}

func init() {
	eenpCmd.Flags().Uint64Var(&heightFlag, "height", uint64(0), "height of the block")
	eenpCmd.MarkFlagRequired("height")
}

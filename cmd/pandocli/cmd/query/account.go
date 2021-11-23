package query

import (
	"encoding/json"
	"fmt"

	"github.com/pandotoken/pando/cmd/pandocli/cmd/utils"
	"github.com/pandotoken/pando/rpc"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	rpcc "github.com/ybbus/jsonrpc"
)

// accountCmd represents the account command.
// Example:
//		pandocli query account --address=0xdf1f3D3eE9430dB3A44aE6B80Eb3E23352BB785E
var accountCmd = &cobra.Command{
	Use:     "account",
	Short:   "Get account status",
	Long:    `Get account status.`,
	Example: `pandocli query account --address=0xdf1f3D3eE9430dB3A44aE6B80Eb3E23352BB785E`,
	Run:     doAccountCmd,
}

func doAccountCmd(cmd *cobra.Command, args []string) {
	client := rpcc.NewRPCClient(viper.GetString(utils.CfgRemoteRPCEndpoint))

	res, err := client.Call("pando.GetAccount", rpc.GetAccountArgs{
		Address: addressFlag, Preview: previewFlag})
	if err != nil {
		utils.Error("Failed to get account details: %v\n", err)
	}
	if res.Error != nil {
		utils.Error("Failed to get account details: %v\n", res.Error)
	}
	json, err := json.MarshalIndent(res.Result, "", "    ")
	if err != nil {
		utils.Error("Failed to parse server response: %v\n%v\n", err, string(json))
	}
	fmt.Println(string(json))
}

func init() {
	accountCmd.Flags().StringVar(&addressFlag, "address", "", "Address of the account")
	accountCmd.Flags().BoolVar(&previewFlag, "preview", false, "Preview account balance from the screened view")
	accountCmd.MarkFlagRequired("address")
}

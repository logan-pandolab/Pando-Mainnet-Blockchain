package rpc

import (
	"encoding/hex"
	"fmt"

	"github.com/pandotoken/pando/common"
	"github.com/pandotoken/pando/ledger/state"
	"github.com/pandotoken/pando/ledger/types"
	"github.com/pandotoken/pando/ledger/vm"
)

// ------------------------------- CallSmartContract -----------------------------------

type CallSmartContractArgs struct {
	SctxBytes string `json:"sctx_bytes"`
}

type CallSmartContractResult struct {
	VmReturn        string            `json:"vm_return"`
	ContractAddress common.Address    `json:"contract_address"`
	GasUsed         common.JSONUint64 `json:"gas_used"`
	VmError         string            `json:"vm_error"`
}

// CallSmartContract calls the smart contract. However, calling a smart contract does NOT modify
// the globally consensus state. It can be used for dry run, or for retrieving info from smart contracts
// without actually spending gas.
func (t *PandoRPCService) CallSmartContract(args *CallSmartContractArgs, result *CallSmartContractResult) (err error) {
	var ledgerState *state.StoreView
	ledgerState, err = t.ledger.GetDeliveredSnapshot()
	if err != nil {
		return err
	}

	blockHeight := ledgerState.Height() + 1 // the view points to the parent of the current block
	if blockHeight < common.HeightEnableSmartContract {
		return fmt.Errorf("Smart contract feature not enabled until block height %v.", common.HeightEnableSmartContract)
	}

	sctxBytes, err := hex.DecodeString(args.SctxBytes)
	if err != nil {
		return err
	}

	tx, err := types.TxFromBytes(sctxBytes)
	if err != nil {
		return fmt.Errorf("Failed to parse SmartContractTx, error: %v", err)
	}
	sctx, ok := tx.(*types.SmartContractTx)
	if !ok {
		return fmt.Errorf("Failed to parse SmartContractTx: %v", args.SctxBytes)
	}

	parentBlock := t.ledger.State().ParentBlock()
	vmRet, contractAddr, gasUsed, vmErr := vm.Execute(parentBlock, sctx, ledgerState)
	ledgerState.Save()

	result.VmReturn = hex.EncodeToString(vmRet)
	result.ContractAddress = contractAddr
	result.GasUsed = common.JSONUint64(gasUsed)
	if vmErr != nil {
		result.VmError = vmErr.Error()
	}

	return nil
}

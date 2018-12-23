package handler

import (
	"reflect"

	"github.com/bluele/hypermint/pkg/abci/types"
	"github.com/bluele/hypermint/pkg/account"
	"github.com/bluele/hypermint/pkg/contract"
	"github.com/bluele/hypermint/pkg/transaction"
)

func NewHandler(am account.AccountMapper, cm *contract.ContractManager) types.Handler {
	return func(ctx types.Context, tx types.Tx) types.Result {
		switch tx := tx.(type) {
		case *transaction.TransferTx:
			return handleTransferTx(ctx, am, tx)
		case *transaction.ContractDeployTx:
			return handleContractDeployTx(ctx, cm, tx)
		case *transaction.ContractCallTx:
			return handleContractCallTx(ctx, cm, tx)
		default:
			errMsg := "Unrecognized Tx type: " + reflect.TypeOf(tx).Name()
			return types.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleTransferTx(ctx types.Context, am account.AccountMapper, tx *transaction.TransferTx) types.Result {
	if err := am.Transfer(ctx, tx.From, tx.Amount, tx.To); err != nil {
		return transaction.ErrFailTransfer(transaction.DefaultCodespace, err.Error()).Result()
	}
	return types.Result{}
}

func handleContractDeployTx(ctx types.Context, cm *contract.ContractManager, tx *transaction.ContractDeployTx) types.Result {
	if _, err := cm.DeployContract(ctx, tx); err != nil {
		return transaction.ErrInvalidDeploy(transaction.DefaultCodespace, err.Error()).Result()
	}
	return types.Result{}
}

func handleContractCallTx(ctx types.Context, cm *contract.ContractManager, tx *transaction.ContractCallTx) types.Result {
	vm, err := cm.GetVM(ctx, tx.Address)
	if err != nil {
		return transaction.ErrInvalidCall(transaction.DefaultCodespace, err.Error()).Result()
	}
	if err := vm.ExecContract(tx.Func); err != nil {
		return transaction.ErrInvalidCall(transaction.DefaultCodespace, err.Error()).Result()
	}
	return types.Result{}
}
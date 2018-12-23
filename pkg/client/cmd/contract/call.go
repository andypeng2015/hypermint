package contract

import (
	"github.com/bluele/hypermint/pkg/client"
	"github.com/bluele/hypermint/pkg/client/helper"
	"github.com/bluele/hypermint/pkg/transaction"
	"github.com/bluele/hypermint/pkg/util"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagContract = "contract"
	flagFunc     = "func"
)

func init() {
	contractCmd.AddCommand(callCmd)
	callCmd.Flags().String(helper.FlagAddress, "", "address")
	callCmd.Flags().String(flagContract, "", "contract address")
	callCmd.Flags().String(flagFunc, "", "function name")
	callCmd.Flags().Uint(flagGas, 0, "gas for tx")
	util.CheckRequiredFlag(callCmd, helper.FlagAddress, flagGas)
}

var callCmd = &cobra.Command{
	Use:   "call",
	Short: "call contract",
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())
		ctx, err := client.NewClientContextFromViper()
		if err != nil {
			return err
		}
		addrs, err := ctx.GetInputAddresses()
		if err != nil {
			return err
		}
		from := addrs[0]
		nonce, err := transaction.GetNonceByAddress(from)
		if err != nil {
			return err
		}
		caddr := common.HexToAddress(viper.GetString(flagContract))
		tx := &transaction.ContractCallTx{
			Address: caddr,
			Func:    viper.GetString(flagFunc),
			Args:    []byte{1},
			CommonTx: transaction.CommonTx{
				From:  from,
				Gas:   uint64(viper.GetInt(flagGas)),
				Nonce: nonce,
			},
		}
		if err := ctx.SignAndBroadcastTx(tx, from); err != nil {
			return err
		}
		return nil
	},
}
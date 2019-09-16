package contract

import (
	"context"
	"fmt"

	"github.com/bluele/hypermint/pkg/client"
	"github.com/bluele/hypermint/pkg/contract/event"
	"github.com/bluele/hypermint/pkg/util"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/types"
)

func init() {
	contractCmd.AddCommand(EventCMD())
}

func EventCMD() *cobra.Command {
	var eventCmd = &cobra.Command{
		Use:   "event",
		Short: "This provides you to pub/sub events",
	}

	// common
	const (
		flagContractAddress = "address"
		flagEventName       = "event.name"
		flagEventValue      = "event.value"
	)

	var subscribeCmd = &cobra.Command{
		Use:   "subscribe",
		Short: "Subscribe Txs using events",
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.BindPFlags(cmd.Flags())
			ctx, err := client.NewClientContextFromViper()
			if err != nil {
				return err
			}
			cl, err := ctx.GetNode()
			if err != nil {
				return err
			}
			if err := cl.OnStart(); err != nil {
				return err
			}
			defer cl.Stop()
			id := common.RandStr(8)
			q := fmt.Sprintf("tm.event='Tx' AND contract.address='%v' AND contract.event.name='%v'", viper.GetString(flagContractAddress), viper.GetString(flagEventName))
			fmt.Printf("subscription-id=%#v query=%#v\n", id, q)
			out, err := cl.Subscribe(context.Background(), id, q)
			if err != nil {
				return err
			}
			for ev := range out {
				etx := ev.Data.(types.EventDataTx)
				fmt.Printf("TxID=0x%x\n", etx.Tx.Hash())
				for _, ev := range etx.Result.Events {
					if ev.Type != "contract" {
						continue
					}
					for _, tag := range ev.Attributes {
						if k := string(tag.GetKey()); k == "event.data" {
							e, err := event.ParseEntry(tag.GetValue())
							if err != nil {
								return err
							}
							fmt.Println(e.String())
						} else if k == "event.name" || k == "address" {
							// skip
						} else {
							fmt.Printf("unknown event: %v\n", tag)
						}
					}
				}
			}
			return nil
		},
	}
	subscribeCmd.Flags().String(flagContractAddress, "", "contract address for subscription")
	subscribeCmd.Flags().String(flagEventName, "", "event name for subscription")
	util.CheckRequiredFlag(subscribeCmd, flagContractAddress, flagEventName)
	eventCmd.AddCommand(subscribeCmd)

	// search
	const (
		flagCount = "count"
	)

	var searchCmd = &cobra.Command{
		Use:   "search",
		Short: "Search Txs using events",
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.BindPFlags(cmd.Flags())
			ctx, err := client.NewClientContextFromViper()
			if err != nil {
				return err
			}
			cl, err := ctx.GetNode()
			if err != nil {
				return err
			}
			contractAddr := ethcmn.HexToAddress(viper.GetString(flagContractAddress))
			q, err := event.MakeEventSearchQuery(
				contractAddr,
				viper.GetString(flagEventName),
				viper.GetString(flagEventValue),
			)
			if err != nil {
				return err
			}
			res, err := cl.TxSearch(q, true, 0, 0)
			if err != nil {
				return err
			}
			if viper.GetBool(flagCount) {
				var count int
				for _, tx := range res.Txs {
					events, err := event.GetEventsByContractAddr(contractAddr, tx)
					if err != nil {
						return err
					}
					if len(events) == 0 {
						continue
					}
					events, err = event.FilterEvents(
						events,
						viper.GetString(flagEventName),
						viper.GetString(flagEventValue),
					)
					if err != nil {
						return err
					}
					count += len(events)
				}
				fmt.Print(count)
				return nil
			} else {
				for _, tx := range res.Txs {
					fmt.Println(tx.TxResult.String())
				}
			}
			return nil
		},
	}

	searchCmd.Flags().String(flagContractAddress, "", "contract address for subscription")
	searchCmd.Flags().String(flagEventName, "", "event name")
	searchCmd.Flags().String(flagEventValue, "", "event value as hex string")
	searchCmd.Flags().Bool(flagCount, false, "if true, only print count of matched txs")
	util.CheckRequiredFlag(searchCmd, flagContractAddress, flagEventName)
	eventCmd.AddCommand(searchCmd)

	return eventCmd
}

package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/kava-labs/kava/x/bep3/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	bep3TxCmd := &cobra.Command{
		Use:   "bep3",
		Short: "bep3 transactions subcommands",
	}

	bep3TxCmd.AddCommand(client.PostCommands(
		GetCmdCreateAtomicSwap(cdc),
		GetCmdClaimAtomicSwap(cdc),
		GetCmdRefundAtomicSwap(cdc),
	)...)

	return bep3TxCmd
}

// GetCmdCreateAtomicSwap cli command for creating atomic swaps
func GetCmdCreateAtomicSwap(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create [to] [recipient-other-chain] [sender-other-chain] [hashed-secret] [timestamp] [coins] [expected-income] [height-span] [cross-chain]",
		Short: "create a new atomic swap",
		Example: fmt.Sprintf("%s tx %s create kava1xy7hrjy9r0algz9w3gzm8u6mrpq97kwta747gj bnb1urfermcg92dwq36572cx4xg84wpk3lfpksr5g7 bnb1uky3me9ggqypmrsvxk7ur6hqkzq7zmv4ed4ng7 bd8f3e7180c3ee5d89c4c6042caf2e8aa4d15782aa924509c8646497d278c902 1583358734 100btc 99btc 360 true --from validator",
			version.ClientName, types.ModuleName),
		Args: cobra.ExactArgs(9),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			from := cliCtx.GetFromAddress() // same as KavaExecutor.DeputyAddress (for cross-chain)
			to, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			recipientOtherChain := args[1] // same as OtherExecutor.DeputyAddress
			senderOtherChain := args[2]

			randomNumberHash, err := types.HexToBytes(args[3])
			if err != nil {
				return err
			}

			timeStamp, err := strconv.ParseInt(args[4], 10, 64)
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoins(args[5])
			if err != nil {
				return err
			}

			expectedIncome := args[6]

			heightSpan, err := strconv.ParseInt(args[7], 10, 64)
			if err != nil {
				return err
			}

			crossChain, err := strconv.ParseBool(args[8])
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateAtomicSwap(
				from, to, recipientOtherChain, senderOtherChain, randomNumberHash,
				timeStamp, coins, expectedIncome, heightSpan, crossChain,
			)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdClaimAtomicSwap cli command for claiming an atomic swap
func GetCmdClaimAtomicSwap(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "claim [swap-id] [random-number]",
		Short:   "claim coins in an atomic swap using the secret number",
		Example: fmt.Sprintf("%s tx %s claim 6682c03cc3856879c8fb98c9733c6b0c30758299138166b6523fe94628b1d3af 56f13e6a5cd397447f8b5f8c82fdb5bbf56127db75269f5cc14e50acd8ac9a4c --from accA", version.ClientName, types.ModuleName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			from := cliCtx.GetFromAddress()

			swapID, err := types.HexToBytes(args[0])
			if err != nil {
				return err
			}

			if len(strings.TrimSpace(args[1])) == 0 {
				return fmt.Errorf("random-number cannot be empty")
			}

			randomNumber, err := types.HexToBytes(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgClaimAtomicSwap(from, swapID, randomNumber)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdRefundAtomicSwap cli command for claiming an atomic swap
func GetCmdRefundAtomicSwap(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "refund [swap-id]",
		Short:   "refund the coins in an atomic swap",
		Example: fmt.Sprintf("%s tx %s refund 6682c03cc3856879c8fb98c9733c6b0c30758299138166b6523fe94628b1d3af --from accA", version.ClientName, types.ModuleName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			from := cliCtx.GetFromAddress()

			swapID, err := types.HexToBytes(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgRefundAtomicSwap(from, swapID)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

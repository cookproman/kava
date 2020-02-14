package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		GetCmdCreateHtlt(cdc),
		GetCmdDepositHtlt(cdc),
	)...)

	return bep3TxCmd
}

// GetCmdCreateHtlt cli command for creating htlts
func GetCmdCreateHtlt(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "create [to] [recipien-other-chain] [hashed-secret] [timestamp] [coins] [expected-income] [height-span] [cross-chain]",
		Short:   "create a new Hashed Time Locked Transaction (HTLT)",
		Example: "bep3 create kava15qdefkmwswysgg4qxgqpqr35k3m49pkx2jdfnw 0x9eB05a790e2De0a047a57a22199D8CccEA6d6D5A 0677bd8a303dd981810f34d8e5cc6507f13b391899b84d3c1be6c6045a17d747 9988776655 100xrp 99xrp 500000 true --from accA",
		Args:    cobra.MinimumNArgs(8),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			from := cliCtx.GetFromAddress()

			// TODO: string -> sdk.AccAddress conversion
			//	- sdk.AccAddressFromBech32(): 'decoding bech32 failed: checksum failed. Expected fefqtn, got 2jdfnw.'
			//	- sdk.AccAddressFromHex(): encoding/hex: invalid byte: U+006B 'k'
			//	- sdk.AccAddress(): creates HTLT with incorrect address length causing HTLT query to fail
			to := from // same as KavaExecutor.DeputyAddress

			if len(args[1]) == 0 {
				return errors.New("recipient-other-chain cannot be empty")
			}
			recipientOtherChain := args[1] // same as OtherExecutor.DeputyAddress

			// TODO: Add as optional arg
			senderOtherChain := ""

			if len(strings.TrimSpace(args[2])) != types.RandomNumberHashLength {
				return fmt.Errorf("random-number-hash should have length %d", types.RandomNumberHashLength)
			}
			hashedSecret := args[2]

			timeStamp, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoins(args[4])
			if err != nil {
				return err
			}
			expectedIncome := args[5]

			heightSpan, err := strconv.ParseInt(args[6], 10, 64)
			if err != nil {
				return err
			}

			crossChain, err := strconv.ParseBool(args[7])
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateHTLT(
				from, to, recipientOtherChain, senderOtherChain, hashedSecret,
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

// GetCmdDepositHtlt cli command for depositing into an htlt
func GetCmdDepositHtlt(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "deposit [swap-id] [coins]",
		Short:   "deposit coins into an existing HTLT",
		Example: "bep3 deposit 6682c03cc3856879c8fb98c9733c6b0c30758299138166b6523fe94628b1d3af 10xrp  --from accA",
		Args:    cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			from := cliCtx.GetFromAddress()

			if len(strings.TrimSpace(args[0])) != types.SwapIDLength {
				return fmt.Errorf("swap-id should have length %d", types.SwapIDLength)
			}
			swapID := args[0]

			if len(strings.TrimSpace(args[1])) == 0 {
				return fmt.Errorf("coins cannot be empty")
			}
			coins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDepositHTLT(from, swapID, coins)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdClaimHtlt cli command for claiming an htlt
func GetCmdClaimHtlt(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "claim [swap-id] [random-number]",
		Short:   "claim coins in an HTLT using the secret random number",
		Example: "bep3 claim 6682c03cc3856879c8fb98c9733c6b0c30758299138166b6523fe94628b1d3af 123456789 --from accA",
		Args:    cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			from := cliCtx.GetFromAddress()

			if len(strings.TrimSpace(args[0])) != types.SwapIDLength {
				return fmt.Errorf("swap-id should have length %d", types.SwapIDLength)
			}
			swapID := args[0]

			if len(strings.TrimSpace(args[1])) == 0 {
				return fmt.Errorf("random-number cannot be empty")
			}
			randomNumber := []byte(args[1])

			msg := types.NewMsgClaimHTLT(from, swapID, randomNumber)

			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdRefundHtlt cli command for claiming an htlt
func GetCmdRefundHtlt(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "refund [swap-id]",
		Short:   "refund the coins in an HTLT",
		Example: "bep3 refund 6682c03cc3856879c8fb98c9733c6b0c30758299138166b6523fe94628b1d3af --from accA",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			from := cliCtx.GetFromAddress()

			if len(strings.TrimSpace(args[0])) != types.SwapIDLength {
				return fmt.Errorf("swap-id should have length %d", types.SwapIDLength)
			}
			swapID := args[0]

			msg := types.NewMsgRefundHTLT(from, swapID)

			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

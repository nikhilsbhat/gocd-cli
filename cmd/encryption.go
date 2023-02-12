package cmd

import (
	"github.com/nikhilsbhat/gocd-cli/pkg/errors"
	"github.com/spf13/cobra"
)

var cipherKey string

func getEncryptionCommand() *cobra.Command {
	encryptionCommand := &cobra.Command{
		Use:   "encryption",
		Short: "Command to encrypt/decrypt plain text value [https://api.gocd.org/current/#encryption]",
		Long: `Command leverages GoCD api [https://api.gocd.org/current/#encryption] during value encryption and 
AES decryption while decrypting encrypted value [https://github.com/nikhilsbhat/gocd-sdk-go/blob/master/encryption.go#L49]`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Usage(); err != nil {
				return err
			}

			return nil
		},
	}
	registerEncryptionFlags(encryptionCommand)
	encryptionCommand.SetUsageTemplate(getUsageTemplate())
	encryptionCommand.AddCommand(getEncryptCommand())
	encryptionCommand.AddCommand(getDecryptCommand())

	return encryptionCommand
}

func getEncryptCommand() *cobra.Command {
	encryptionCommand := &cobra.Command{
		Use:     "encrypt",
		Short:   "Command to encrypt plain text value [https://api.gocd.org/current/#encryption]",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: setGoCDClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return &errors.AuthError{Message: "args cannot be more than one, only one encrypted/plain text must be passed"}
			}

			response, err := client.EncryptText(args[0])
			if err != nil {
				return err
			}

			if err = render(response); err != nil {
				return err
			}

			return nil
		},
	}

	encryptionCommand.SetUsageTemplate(getUsageTemplate())

	return encryptionCommand
}

func getDecryptCommand() *cobra.Command {
	encryptionCommand := &cobra.Command{
		Use:     "decrypt",
		Short:   "Command to decrypt encrypted value [https://github.com/nikhilsbhat/gocd-sdk-go/blob/master/encryption.go#L49]",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: setGoCDClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return &errors.AuthError{Message: "args cannot be more than one, only one encrypted/plain text must be passed"}
			}

			response, err := client.DecryptText(args[0], cipherKey)
			if err != nil {
				return err
			}

			if err = render(response); err != nil {
				return err
			}

			return nil
		},
	}

	encryptionCommand.SetUsageTemplate(getUsageTemplate())

	return encryptionCommand
}

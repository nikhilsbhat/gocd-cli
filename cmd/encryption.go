package cmd

import (
	"github.com/nikhilsbhat/gocd-cli/pkg/errors"
	"github.com/spf13/cobra"
	"os"
)

var (
	cipherKey     string
	cipherKeyPath string
)

func registerEncryptionCommand() *cobra.Command {
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

	encryptionCommand.SetUsageTemplate(getUsageTemplate())
	encryptionCommand.AddCommand(getEncryptCommand())
	encryptionCommand.AddCommand(getDecryptCommand())

	for _, command := range encryptionCommand.Commands() {
		command.SilenceUsage = true
	}

	return encryptionCommand
}

func getEncryptCommand() *cobra.Command {
	encryptionCommand := &cobra.Command{
		Use:     "encrypt",
		Short:   "Command to encrypt plain text value [https://api.gocd.org/current/#encryption]",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return &errors.MoreArgError{Message: "encrypted/plain"}
			}

			response, err := client.EncryptText(args[0])
			if err != nil {
				return err
			}

			if err = cliRenderer.Render(response); err != nil {
				return err
			}

			return nil
		},
	}

	encryptionCommand.SetUsageTemplate(getUsageTemplate())

	return encryptionCommand
}

func getDecryptCommand() *cobra.Command {
	decryptionCommand := &cobra.Command{
		Use:     "decrypt",
		Short:   "Command to decrypt encrypted value [https://github.com/nikhilsbhat/gocd-sdk-go/blob/master/encryption.go#L49]",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return &errors.MoreArgError{Message: "encrypted/plain"}
			}

			if len(cipherKeyPath) != 0 {
				out, err := os.ReadFile(cipherKeyPath)
				if err != nil {
					return err
				}

				cipherKey = string(out)
			}

			response, err := client.DecryptText(args[0], cipherKey)
			if err != nil {
				return err
			}

			if err = cliRenderer.Render(response); err != nil {
				return err
			}

			return nil
		},
	}

	registerEncryptionFlags(decryptionCommand)

	decryptionCommand.SetUsageTemplate(getUsageTemplate())

	return decryptionCommand
}

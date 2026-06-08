package cmd

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/nikhilsbhat/common/content"
	"github.com/nikhilsbhat/gocd-cli/pkg/query"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/spf13/cobra"
)

func registerTemplatesCommand() *cobra.Command {
	templatesCommand := &cobra.Command{
		Use:   "template",
		Short: "Command to operate on templates present in GoCD [https://api.gocd.org/current/#pipeline-template-config]",
		Long: `Command leverages GoCD template config apis' [https://api.gocd.org/current/#pipeline-template-config] to
GET/CREATE/UPDATE/DELETE and list GoCD templates`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Usage()
		},
	}

	templatesCommand.SetUsageTemplate(getUsageTemplate())

	templatesCommand.AddCommand(getTemplatesCommand())
	templatesCommand.AddCommand(getTemplateCommand())
	templatesCommand.AddCommand(createTemplateCommand())
	templatesCommand.AddCommand(updateTemplateCommand())
	templatesCommand.AddCommand(deleteTemplateCommand())
	templatesCommand.AddCommand(listTemplatesCommand())
	templatesCommand.AddCommand(getTemplateParametersCommand())
	templatesCommand.AddCommand(getTemplateAuthorizationCommand())
	templatesCommand.AddCommand(updateTemplateAuthorizationCommand())

	for _, command := range templatesCommand.Commands() {
		command.SilenceUsage = true
	}

	return templatesCommand
}

func getTemplatesCommand() *cobra.Command {
	getTemplatesCmd := &cobra.Command{
		Use:     "get-all",
		Short:   "Command to GET all templates present in GoCD [https://api.gocd.org/current/#get-all-templates]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		Example: `gocd-cli template get-all --query "templates.[*] | name"`,
		RunE: func(_ *cobra.Command, _ []string) error {
			for {
				response, err := client.GetTemplates()
				if err != nil {
					return err
				}

				if len(jsonQuery) != 0 {
					cliLogger.Debugf(queryEnabledMessage, jsonQuery)

					baseQuery, err := query.SetQuery(response, jsonQuery)
					if err != nil {
						return err
					}

					cliLogger.Debug(baseQuery.Print())

					return cliRenderer.Render(baseQuery.RunQuery())
				}

				if err = cliRenderer.Render(response); err != nil {
					return err
				}

				if !cliCfg.Watch {
					break
				}

				time.Sleep(cliCfg.WatchInterval)
			}

			return nil
		},
	}

	return getTemplatesCmd
}

func getTemplateCommand() *cobra.Command {
	getTemplateCmd := &cobra.Command{
		Use:     "get",
		Short:   "Command to GET a specific template present in GoCD [https://api.gocd.org/current/#get-a-template]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		Example: `gocd-cli template get build-template --query "stages.[*] | name" -o yaml`,
		RunE: func(_ *cobra.Command, args []string) error {
			for {
				response, err := client.GetTemplate(args[0])
				if err != nil {
					return err
				}

				if len(jsonQuery) != 0 {
					cliLogger.Debugf(queryEnabledMessage, jsonQuery)

					baseQuery, err := query.SetQuery(response, jsonQuery)
					if err != nil {
						return err
					}

					cliLogger.Debug(baseQuery.Print())

					return cliRenderer.Render(baseQuery.RunQuery())
				}

				if err = cliRenderer.Render(response); err != nil {
					return err
				}

				if !cliCfg.Watch {
					break
				}

				time.Sleep(cliCfg.WatchInterval)
			}

			return nil
		},
	}

	return getTemplateCmd
}

func createTemplateCommand() *cobra.Command {
	createTemplateCmd := &cobra.Command{
		Use:     "create",
		Short:   "Command to CREATE a template with all specified configuration [https://api.gocd.org/current/#create-a-template]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		Example: `gocd-cli template create --from-file template.yaml --log-level debug`,
		RunE:    createTemplate,
	}

	return createTemplateCmd
}

func updateTemplateCommand() *cobra.Command {
	updateTemplateCmd := &cobra.Command{
		Use:     "update",
		Short:   "Command to UPDATE the template with the latest specified configuration [https://api.gocd.org/current/#update-a-template]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		Example: `gocd-cli template update --from-file template.yaml --log-level debug`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			template, object, err := readTemplate(cmd)
			if err != nil {
				return err
			}

			templateFetched, err := client.GetTemplate(template.Name)
			if err != nil && !strings.Contains(err.Error(), "404") {
				return err
			}

			if create {
				if reflect.DeepEqual(templateFetched, gocd.Template{}) {
					return createTemplate(cmd, nil)
				}
			}

			if len(template.ETAG) == 0 {
				template.ETAG = templateFetched.ETAG
			}

			cliShellReadConfig.ShellMessage = fmt.Sprintf(updateMessage, "template", template.Name)

			existing, err := diffCfg.String(templateFetched)
			if err != nil {
				return err
			}

			if err = cliCfg.CheckDiffAndAllow(existing, object.String()); err != nil {
				return err
			}

			updatedTemplate, err := client.UpdateTemplate(template)
			if err != nil {
				return err
			}

			if err = cliRenderer.Render(fmt.Sprintf("template %s updated successfully", template.Name)); err != nil {
				return err
			}

			return cliRenderer.Render(updatedTemplate)
		},
	}

	updateTemplateCmd.PersistentFlags().BoolVarP(&create, "create", "", false,
		"if a template by this name doesn't already exist, run create")

	return updateTemplateCmd
}

func deleteTemplateCommand() *cobra.Command {
	deleteTemplateCmd := &cobra.Command{
		Use:     "delete",
		Short:   "Command to DELETE the specified template from GoCD [https://api.gocd.org/current/#delete-a-template]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		Example: `gocd-cli template delete build-template`,
		RunE: func(_ *cobra.Command, args []string) error {
			templateName := args[0]
			cliShellReadConfig.ShellMessage = fmt.Sprintf("do you want to delete template '%s' [y/n]", templateName)

			if !cliCfg.Yes {
				contains, option := cliShellReadConfig.Reader()
				if !contains {
					cliLogger.Fatalln(inputValidationFailureMessage)
				}

				if option.Short == "n" {
					cliLogger.Warn(optingOutMessage)

					os.Exit(0)
				}
			}

			if err := client.DeleteTemplate(templateName); err != nil {
				return err
			}

			return cliRenderer.Render(fmt.Sprintf("template deleted: %s", templateName))
		},
	}

	return deleteTemplateCmd
}

func listTemplatesCommand() *cobra.Command {
	listTemplatesCmd := &cobra.Command{
		Use:     "list",
		Short:   "Command to LIST all templates present in GoCD [https://api.gocd.org/current/#get-all-templates]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		Example: `gocd-cli template list`,
		RunE: func(_ *cobra.Command, _ []string) error {
			for {
				response, err := client.GetTemplates()
				if err != nil {
					return err
				}

				var templates []string
				for _, template := range response.Templates {
					templates = append(templates, template.Name)
				}

				if err = cliRenderer.Render(strings.Join(templates, "\n")); err != nil {
					return err
				}

				if !cliCfg.Watch {
					break
				}

				time.Sleep(cliCfg.WatchInterval)
			}

			return nil
		},
	}

	return listTemplatesCmd
}

func getTemplateParametersCommand() *cobra.Command {
	getTemplateParametersCmd := &cobra.Command{
		Use:     "parameters",
		Short:   "Command to GET parameters for a specific template present in GoCD [https://api.gocd.org/current/#get-template-parameters]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		Example: `gocd-cli template parameters build-template`,
		RunE: func(_ *cobra.Command, args []string) error {
			response, err := client.GetTemplateParameters(args[0])
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	return getTemplateParametersCmd
}

func getTemplateAuthorizationCommand() *cobra.Command {
	getTemplateAuthorizationCmd := &cobra.Command{
		Use:     "get-authorization",
		Short:   "Command to GET authorization for a specific template present in GoCD [https://api.gocd.org/current/#get-template-authorization]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		Example: `gocd-cli template get-authorization build-template`,
		RunE: func(_ *cobra.Command, args []string) error {
			response, err := client.GetTemplateAuthorization(args[0])
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	return getTemplateAuthorizationCmd
}

func updateTemplateAuthorizationCommand() *cobra.Command {
	updateTemplateAuthorizationCmd := &cobra.Command{
		Use:     "update-authorization",
		Short:   "Command to UPDATE authorization for a specific template present in GoCD [https://api.gocd.org/current/#update-template-authorization]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		Example: `gocd-cli template update-authorization build-template --from-file template-authorization.yaml`,
		RunE: func(cmd *cobra.Command, args []string) error {
			authorization, object, err := readTemplateAuthorization(cmd)
			if err != nil {
				return err
			}

			authorizationFetched, err := client.GetTemplateAuthorization(args[0])
			if err != nil {
				return err
			}

			if len(authorization.ETAG) == 0 {
				authorization.ETAG = authorizationFetched.ETAG
			}

			cliShellReadConfig.ShellMessage = fmt.Sprintf(updateMessage, "template-authorization", args[0])

			existing, err := diffCfg.String(authorizationFetched)
			if err != nil {
				return err
			}

			if err = cliCfg.CheckDiffAndAllow(existing, object.String()); err != nil {
				return err
			}

			updatedAuthorization, err := client.UpdateTemplateAuthorization(args[0], authorization)
			if err != nil {
				return err
			}

			if err = cliRenderer.Render(fmt.Sprintf("template authorization %s updated successfully", args[0])); err != nil {
				return err
			}

			return cliRenderer.Render(updatedAuthorization)
		},
	}

	return updateTemplateAuthorizationCmd
}

func createTemplate(cmd *cobra.Command, _ []string) error {
	template, _, err := readTemplate(cmd)
	if err != nil {
		return err
	}

	createdTemplate, err := client.CreateTemplate(template)
	if err != nil {
		return err
	}

	if err = cliRenderer.Render(fmt.Sprintf("template %s created successfully", template.Name)); err != nil {
		return err
	}

	return cliRenderer.Render(createdTemplate)
}

func readTemplate(cmd *cobra.Command) (gocd.Template, content.Object, error) {
	return readConfigObject[gocd.Template](cmd)
}

func readTemplateAuthorization(cmd *cobra.Command) (gocd.TemplateAuthorization, content.Object, error) {
	return readConfigObject[gocd.TemplateAuthorization](cmd)
}

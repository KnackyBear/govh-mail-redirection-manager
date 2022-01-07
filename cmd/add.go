package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var cmdAdd = &cobra.Command{
	Use:               "add",
	Short:             "Add a new mail redirection",
	ValidArgsFunction: cobra.NoFileCompletions,
	RunE: func(cmd *cobra.Command, args []string) error {
		type redirectionBody struct {
			From      string `json:"from"`
			To        string `json:"to"`
			LocalCopy bool   `json:"localCopy"`
		}

		type redirectionResponse struct {
			Account string `json:"account"`
			Id      string `json:"id"`
			Domain  string `json:"domain"`
			Date    string `json:"date"`
			Action  string `json:"action"`
			Type    string `json:"type"`
		}

		payload := redirectionBody{
			From: fromFlag,
			To:   toFlag,
		}

		response := redirectionResponse{}

		var domain string
		if Domain == "" {
			email := strings.Split(fromFlag, "@")
			if len(email) == 2 {
				domain = email[1]
			} else {
				return fmt.Errorf("Domain not found in %v", fromFlag)
			}
		} else {
			domain = Domain
		}

		if err := OvhClient.Post(fmt.Sprintf("/email/domain/%s/redirection", domain), payload, &response); err != nil {
			return fmt.Errorf("error: %q", err)
		}

		fmt.Printf("New redirection from '%s' to '%s' added.", payload.From, payload.To)
		return nil
	},
}

func init() {
	cmdAdd.Flags().StringVar(&fromFlag, "from", "", "Email of redirection")
	cmdAdd.Flags().StringVar(&toFlag, "to", "", "Email of destination")
	cmdAdd.MarkFlagRequired("from")
	cmdAdd.MarkFlagRequired("to")
	cmdAdd.RegisterFlagCompletionFunc("from", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	})
	cmdAdd.RegisterFlagCompletionFunc("to", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	})

	RootCmd.AddCommand(cmdAdd)
}

package cmd

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

var cmdRemove = &cobra.Command{
	Use:               "remove",
	Short:             "Remove a mail redirection",
	ValidArgsFunction: cobra.NoFileCompletions,
	RunE: func(cmd *cobra.Command, args []string) error {
		ids := []string{}

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

		if err := OvhClient.Get(fmt.Sprintf("/email/domain/%s/redirection?from=%s", domain, url.QueryEscape(fromFlag)), &ids); err != nil {
			return fmt.Errorf("error: %q", err)
		}

		type redirectionResponse struct {
			Account string `json:"account"`
			Id      string `json:"id"`
			Domain  string `json:"domain"`
			Date    string `json:"date"`
			Action  string `json:"action"`
			Type    string `json:"type"`
		}

		for _, id := range ids {
			response := redirectionResponse{}
			if err := OvhClient.Delete(fmt.Sprintf("/email/domain/%s/redirection/%s", domain, id), &response); err != nil {
				fmt.Printf("Can't delete redirection %s !", id)
			} else {
				fmt.Printf("Redirection %s deleted.", id)
			}
		}
		return nil
	},
}

func init() {
	cmdRemove.Flags().StringVar(&fromFlag, "from", "", "Email of redirection")
	cmdRemove.MarkFlagRequired("from")
	cmdRemove.RegisterFlagCompletionFunc("from", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	})
	RootCmd.AddCommand(cmdRemove)
}

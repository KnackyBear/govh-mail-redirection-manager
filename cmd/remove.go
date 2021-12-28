package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

var cmdRemove = &cobra.Command{
	Use:               "remove",
	Short:             "Remove a mail redirection",
	ValidArgsFunction: cobra.NoFileCompletions,
	Run: func(cmd *cobra.Command, args []string) {
		ids := []string{}
		if err := OvhClient.Get(fmt.Sprintf("/email/domain/%s/redirection?from=%s", Domain, url.QueryEscape(fromFlag)), &ids); err != nil {
			fmt.Printf("Error: %q\n", err)
			return
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
			if err := OvhClient.Delete(fmt.Sprintf("/email/domain/%s/redirection/%s", Domain, id), &response); err != nil {
				fmt.Printf("Can't delete redirection %s !", id)
			} else {
				fmt.Printf("Redirection %s deleted.", id)
			}
		}
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

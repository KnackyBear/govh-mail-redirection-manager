package cmd

import (
	"fmt"
	"net/url"

	"github.com/julienvinet/govh-mrm/pkg/utils"
	"github.com/spf13/cobra"
)

var cmdRemove = &cobra.Command{
	Use:               "remove <redirection mail>",
	Short:             "Remove a mail redirection",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: cobra.NoFileCompletions,
	RunE: func(cmd *cobra.Command, args []string) error {
		ids := []string{}

		var domain string
		var err error
		if Domain == "" {
			if domain, err = utils.GetDomain(args[0]); err != nil {
				return err
			}
			if !utils.StrArrayContains(currentConfig.Domain, domain) {
				return fmt.Errorf("Unknown domain %v", domain)
			}
		} else {
			domain = Domain
		}

		if err := OvhClient.Get(fmt.Sprintf("/email/domain/%s/redirection?from=%s", domain, url.QueryEscape(args[0])), &ids); err != nil {
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
				fmt.Printf("Can't delete redirection %s (%s) !", args[0], id)
			} else {
				fmt.Printf("Redirection %s (%s) deleted.", args[0], id)
			}
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(cmdRemove)
}

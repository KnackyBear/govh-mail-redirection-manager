package cmd

import (
	"fmt"

	"github.com/julienvinet/govh-mrm/pkg/utils"
	"github.com/spf13/cobra"
)

var cmdAdd = &cobra.Command{
	Use:               "add <redirection mail> <destination mail>",
	Short:             "Add a new mail redirection",
	Args:              cobra.ExactArgs(2),
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
			From: args[0],
			To:   args[1],
		}

		response := redirectionResponse{}

		var domain string
		var err error

		if Domain == "" {
			if domain, err = utils.GetDomain(payload.From); err != nil {
				return err
			}
			if !utils.StrArrayContains(currentConfig.Domain, domain) {
				return fmt.Errorf("Unknown domain %v", domain)
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
	RootCmd.AddCommand(cmdAdd)
}

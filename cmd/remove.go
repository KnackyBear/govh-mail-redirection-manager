package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

var cmdRemove = &cobra.Command{
	Use:   "remove",
	Short: "Remove a mail redirection",
	Long:  `This command remove a mail redirection using OVH Provider.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if strings.TrimSpace(OptFrom) == "" {
			return errors.New("missing argument From")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ids := []string{}
		if err := OvhClient.Get(fmt.Sprintf("/email/domain/%s/redirection?from=%s&to=%s", Domain, url.QueryEscape(OptFrom), url.QueryEscape(OptTo)), &ids); err != nil {
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
	RootCmd.AddCommand(cmdRemove)
}

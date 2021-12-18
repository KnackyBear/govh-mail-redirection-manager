package cmd

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

var fromFlag string

var cmdRemove = &cobra.Command{
	Use:       "remove",
	Short:     "Remove a mail redirection",
	ValidArgs: []string{"--from="},
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
	cmdRemove.PersistentFlags().StringVar(&fromFlag, "from", "", "Name of redirection")
	cmdRemove.MarkFlagRequired("from")
	cmdRemove.RegisterFlagCompletionFunc("from", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		redirections, _ := ListRedirection("", "")
		strArray := make([]string, len(redirections))
		for i, arg := range redirections {
			strArray[i] = arg.From
		}

		prefixArray := make([]string, 0)
		for _, str := range strArray {
			if strings.HasPrefix(str, toComplete) {
				prefixArray = append(prefixArray, str)
			}
		}

		if len(prefixArray) > 0 {
			return prefixArray, cobra.ShellCompDirectiveNoFileComp
		} else {
			return strArray, cobra.ShellCompDirectiveNoFileComp
		}
	})

	RootCmd.AddCommand(cmdRemove)
}

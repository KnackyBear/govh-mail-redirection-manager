package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

type RedirectionAccess struct {
	From string `json:"from"`
	To   string `json:"to"`
}

var cmdList = &cobra.Command{
	Use:   "list",
	Short: "List mail redirection",
	Long: `This command list all mail redirection from OVH. 
You can filter the search with the use of flags 'from' and 'to'.`,
	ValidArgs: []string{"from", "to"},
	Run: func(cmd *cobra.Command, args []string) {
		redirections, err := ListRedirection(OptFrom, OptTo)

		if err != nil {
			fmt.Printf("Error: %q\n", err)
			return
		} else {
			for _, redirection := range redirections {
				fmt.Printf("%-50s %s\n", redirection.From, redirection.To)
			}
		}
	},
}

func ListRedirection(from string, to string) ([]RedirectionAccess, error) {
	redirections := make([]RedirectionAccess, 0)
	ids := []string{}
	if err := OvhClient.Get(fmt.Sprintf("/email/domain/%s/redirection?from=%s&to=%s", Domain, url.QueryEscape(from), url.QueryEscape(to)), &ids); err != nil {
		return nil, err
	}

	for _, id := range ids {
		redirection := RedirectionAccess{}
		OvhClient.Get(fmt.Sprintf("/email/domain/%s/redirection/%s", Domain, id), &redirection)
		redirections = append(redirections, redirection)
	}
	return redirections, nil
}

func init() {
	RootCmd.AddCommand(cmdList)
}

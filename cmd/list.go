package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

var cmdList = &cobra.Command{
	Use:   "list",
	Short: "List mail redirection",
	Long: `This command list all mail redirection from OVH. 
You can filter the search with the use of flags 'from' and 'to'.`,
	Run: func(cmd *cobra.Command, args []string) {
		ids := []string{}
		if err := OvhClient.Get(fmt.Sprintf("/email/domain/%s/redirection?from=%s&to=%s", Domain, url.QueryEscape(OptFrom), url.QueryEscape(OptTo)), &ids); err != nil {
			fmt.Printf("Error: %q\n", err)
			return
		}

		type redirectionAccess struct {
			From string `json:"from"`
			To   string `json:"to"`
		}

		for _, id := range ids {
			redirection := redirectionAccess{}
			OvhClient.Get(fmt.Sprintf("/email/domain/%s/redirection/%s", Domain, id), &redirection)
			fmt.Printf("%-50s %s\n", redirection.From, redirection.To)
		}
	},
}

func init() {
	RootCmd.AddCommand(cmdList)
}

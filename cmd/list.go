package cmd

import (
	"fmt"
	"net/url"
	"time"

	"github.com/julienvinet/govh-mrm/pkg/utils"
	"github.com/spf13/cobra"
)

var bar = utils.Bar{}
var cpt = 0

type RedirectionAccess struct {
	From string `json:"from"`
	To   string `json:"to"`
}

var cmdList = &cobra.Command{
	Use:   "list",
	Short: "List mail redirection",
	Long: `This command list all mail redirection from OVH. 
You can filter the search with the use of flags 'from' and 'to'.`,
	DisableFlagsInUseLine: true,
	ValidArgsFunction:     cobra.NoFileCompletions,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		redirections, err := ListRedirection(fromFlag, toFlag)
		if err != nil {
			fmt.Printf("Error: %q\n", err)
			return
		} else {
			for _, redirection := range redirections {
				fmt.Printf("%-50s %s\n", redirection.From, redirection.To)
			}
		}
		fmt.Printf("\nRequest took %v\n", time.Since(start).Round(time.Millisecond))
	},
}

func getRedirectionDetails(id string, rc chan *RedirectionAccess) {
	redirection := RedirectionAccess{}
	OvhClient.Get(fmt.Sprintf("/email/domain/%s/redirection/%s", Domain, id), &redirection)
	cpt = cpt + 1
	bar.Play(int64(cpt))
	rc <- &redirection
}

func ListRedirection(from string, to string) ([]RedirectionAccess, error) {
	redirectionChan := make(chan *RedirectionAccess)
	redirections := make([]RedirectionAccess, 0)
	ids := []string{}
	if err := OvhClient.Get(fmt.Sprintf("/email/domain/%s/redirection?from=%s&to=%s", Domain, url.QueryEscape(from), url.QueryEscape(to)), &ids); err != nil {
		return nil, err
	}

	bar.NewOption(0, int64(len(ids)))

	for _, id := range ids {
		go getRedirectionDetails(id, redirectionChan)
		redirection := <-redirectionChan
		redirections = append(redirections, *redirection)
	}
	bar.Finish()

	return redirections, nil
}

func init() {
	cmdList.Flags().StringVar(&fromFlag, "from", "", "Email of redirection")
	cmdList.Flags().StringVar(&toFlag, "to", "", "Email of destination")
	cmdList.RegisterFlagCompletionFunc("from", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	})
	cmdList.RegisterFlagCompletionFunc("to", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	})
	RootCmd.AddCommand(cmdList)
}

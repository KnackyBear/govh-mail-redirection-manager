package cmd

import (
	"fmt"
	"net/url"
	"time"

	"github.com/julienvinet/govh-mrm/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	bar = utils.Bar{}
	cpt = 0
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
	DisableFlagsInUseLine: true,
	ValidArgsFunction:     cobra.NoFileCompletions,
	RunE: func(cmd *cobra.Command, args []string) error {
		domains := make([]string, 0)
		if Domain == "" {
			domains = currentConfig.Domain
		} else {
			domains = append(domains, Domain)
		}
		start := time.Now()
		redirections, err := ListRedirection(fromFlag, toFlag, domains)
		if err != nil {
			return fmt.Errorf("Error: %q\n", err)
		} else {
			for _, redirection := range redirections {
				fmt.Printf("%-50s %s\n", redirection.From, redirection.To)
			}
		}
		fmt.Printf("\nRequest took %v\n", time.Since(start).Round(time.Millisecond))
		return nil
	},
}

func getRedirectionDetails(id string, domain string, rc chan *RedirectionAccess) {
	redirection := RedirectionAccess{}
	//fmt.Printf("> /email/domain/%s/redirection/%s\n", domain, id)
	if err := OvhClient.Get(fmt.Sprintf("/email/domain/%s/redirection/%s", domain, id), &redirection); err == nil {
		rc <- &redirection
	}
	cpt = cpt + 1
	bar.Play(int64(cpt))
}

func ListRedirection(from string, to string, domains []string) ([]RedirectionAccess, error) {
	redirectionChan := make(chan *RedirectionAccess)
	redirections := make([]RedirectionAccess, 0)

	ids := make(map[string]string, 0)
	for _, domain := range domains {
		cids := make([]string, 0)
		if err := OvhClient.Get(fmt.Sprintf("/email/domain/%s/redirection?from=%s&to=%s", domain, url.QueryEscape(from), url.QueryEscape(to)), &cids); err != nil {
			return nil, err
		}
		for _, cid := range cids {
			ids[cid] = domain
		}
	}

	bar.NewOption(0, int64(len(ids)))

	for id, dom := range ids {
		go getRedirectionDetails(id, dom, redirectionChan)
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

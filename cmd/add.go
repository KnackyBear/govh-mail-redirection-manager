package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmdAdd = &cobra.Command{
	Use:       "add",
	Short:     "Add a new mail redirection",
	ValidArgs: []string{"from", "to"},
	Args:      cobra.ExactValidArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
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
			From: OptFrom,
			To:   OptTo,
		}

		response := redirectionResponse{}

		if err := OvhClient.Post(fmt.Sprintf("/email/domain/%s/redirection", Domain), payload, &response); err != nil {
			fmt.Printf("Error: %q\n", err)
			return
		}

		fmt.Printf("New redirection from '%s' to '%s' added.", payload.From, payload.To)
	},
}

func init() {
	cmdAdd.MarkFlagRequired("from")
	cmdAdd.MarkFlagRequired("to")

	RootCmd.AddCommand(cmdAdd)
}

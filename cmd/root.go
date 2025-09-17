package cmd

import (
	fn "P2PMail/internal"
	"os"

	"github.com/spf13/cobra"
)

//trying bu using Cobra
var root = &cobra.Command{
	Use: "p2p",
	Long: `/**
 *     ███████████   ████████  ███████████    
 *    ▒▒███▒▒▒▒▒███ ███▒▒▒▒███▒▒███▒▒▒▒▒███   
 *     ▒███    ▒███▒▒▒    ▒███ ▒███    ▒███   
 *     ▒██████████    ███████  ▒██████████    
 *     ▒███▒▒▒▒▒▒    ███▒▒▒▒   ▒███▒▒▒▒▒▒     
 *     ▒███         ███      █ ▒███           
 *     █████       ▒██████████ █████          
 *    ▒▒▒▒▒        ▒▒▒▒▒▒▒▒▒▒ ▒▒▒▒▒           
 *                                            
 *                                            
 *                                            
 */` + "\nA CLI tool for splitting, encrypting, and transmitting files via email with MIME formatting and Gmail integration.",
	Short: "A CLI tool for splitting, encrypting, and transmitting files via email with MIME formatting and Gmail integration.",
	Aliases: []string{"mail", "ptop"},
	Args: cobra.ExactArgs(1),
}

var add = &cobra.Command{
	Use: "add [filename]",
	Short: "Create a Encryted chucks of the given file and store it in MetaData",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := fn.AddFile(args[0])
		if err != nil {
			fn.ErrPrinter(err)
			os.Exit(1)
		}
	},
	//DisableFlagParsing: true,
}

//Yuvaraj work on Push command
var push = &cobra.Command{
	Use: "push [id] [to]",
	Short: "To push the file to others via mail (any mail servies)",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := fn.IsValidMail(args[1]); if err != nil {
			fn.ErrPrinter(err)
			os.Exit(1)
		}
		err = fn.PushFile(args[0], args[1])
		if err != nil {
			fn.ErrPrinter(err)
			os.Exit(1)
		}
	},
}

var pull = &cobra.Command{

}

var reset = &cobra.Command{
	Use: "reset",
	Short: "This will clear the all data in MetaData.json (which is database of your history)",
	SuggestFor: []string{"p2p reset"},
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		err := fn.ClearMetaDataFile("MetaData.json")
		if err != nil {
			fn.ErrPrinter(err)
			os.Exit(1)
		}
	},
}

func Exe() {
	root.AddCommand(add)
	root.AddCommand(reset)
	root.AddCommand(push)
	root.AddCommand(pull)
	err := root.Execute()
	if err != nil {
		fn.ErrPrinter(err)
	}
}

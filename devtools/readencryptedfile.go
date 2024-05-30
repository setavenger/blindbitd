package main

import (
	"fmt"
	"log"
	"os"

	"github.com/setavenger/blindbitd/src/database"
	"github.com/spf13/cobra"
)

var (
	path string
	pass string

	rencfileCmd = &cobra.Command{
		Use:   "rencfile",
		Short: "Reads files encrypted by blindbitd",
		Long:  `Provide a path to the file and a corresponding password to decrypt the file`,
		Run: func(cmd *cobra.Command, args []string) {
			encryptedData, err := os.ReadFile(path)
			if err != nil {
				log.Fatalln(err)
			}

			decryptedData, err := database.DecryptWithPass(encryptedData, []byte(pass))
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(string(decryptedData))
		},
	}
)

func init() {
	RootCmd.AddCommand(rencfileCmd)

	rencfileCmd.PersistentFlags().StringVar(&path, "path", "", "sets the path to the file which should be unlocked")
	rencfileCmd.PersistentFlags().StringVar(&pass, "pass", "", "set the password which should unlock the file")
}

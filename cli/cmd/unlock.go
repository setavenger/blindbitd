package cmd

import (
	"context"
	"fmt"
	"github.com/setavenger/blindbitd/cli/lib"
	"github.com/setavenger/blindbitd/pb"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"google.golang.org/grpc"
	"log"
	"os"
)

// unlockCmd represents the unlock command
var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Unlocks the daemon",
	Long: `This command has to be used before most commands can be used. The daemon is locked on startup with the encryption password. 
The encryption password was set during wallet creation.  
`,
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := lib.NewClient(socketPath)
		defer func(conn *grpc.ClientConn) {
			err := conn.Close()
			if err != nil {
				panic(err)
			}
		}(conn)

		fmt.Print("Enter password: ")
		password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Println("Error reading password")
			return
		}
		fmt.Println()

		response, err := client.Unlock(context.Background(), &pb.PasswordRequest{Password: string(password)})
		if err != nil {
			log.Fatalf("%v", err)
		}
		if response.Success {
			fmt.Println("unlock successfully")
		} else {
			fmt.Println("unlock failed")
			fmt.Println(response.Error)
		}
	},
}

func init() {
	RootCmd.AddCommand(unlockCmd)
}

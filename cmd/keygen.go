package cmd

import (
	"fmt"
	"log"
	"s3-mini/internal/security"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)


var role string

var keygenCmd = &cobra.Command{
	Use: "keygen",
	Short: "Generate an API Key",
	Run: func(cmd *cobra.Command, args []string) {
		ks, _ := security.NewKeyStore("./keys/access_keys.json")

		newKey := "sk_" + uuid.New().String()

		if err := ks.CreateKey(newKey, role); err != nil {
			log.Fatalf("failed to save key: %v", err)
		}

		fmt.Printf("created new key:\n")

		fmt.Printf("	Key:	%s\n",newKey)
		fmt.Printf("	Role:	%s\n",role)	

	},
}


func init() {
	rootCmd.AddCommand(keygenCmd)
	keygenCmd.Flags().StringVar(&role, "role", "WRITE", "Permission Level (READ, WRITE, ADMIN)")
}
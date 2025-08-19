package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"go.starlark.net/starlark"

	"tidbyt.dev/pixlet/runtime"
)

const PublicKeysetJSON = `{"primaryKeyId":3245420378,"key":[{"keyData":{"typeUrl":"type.googleapis.com/google.crypto.tink.HpkePublicKey","value":"EgYIARABGAIaIBypDO7WQjqksjF8V6Yb6H5+ELre3YVX28wKAuPjiFtJ","keyMaterialType":"ASYMMETRIC_PUBLIC"},"status":"ENABLED","keyId":3245420378,"outputPrefixType":"TINK"}]}`

var EncryptCmd = &cobra.Command{
	Use:     "encrypt [app ID] [secret value]...",
	Short:   "Encrypt a secret for use in the Tidbyt community repo",
	Example: "encrypt weather my-top-secretweather-api-key-123456",
	Args:    cobra.MinimumNArgs(2),
	Run:     encrypt,
}

func encrypt(cmd *cobra.Command, args []string) {
	sek := &runtime.SecretEncryptionKey{
		PublicKeysetJSON: []byte(PublicKeysetJSON),
	}

	appID := args[0]
	encrypted := make([]string, len(args)-1)

	for i, val := range args[1:] {
		var err error
		encrypted[i], err = sek.Encrypt(appID, val)
		if err != nil {
			log.Fatalf("encrypting value: %v", err)
		}
	}

	for _, val := range encrypted {
		fmt.Println(starlark.String(val).String())
	}
}

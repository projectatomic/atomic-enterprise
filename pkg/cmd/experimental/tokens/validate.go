package tokens

import (
	"encoding/json"
	"fmt"

	kclient "github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/golang/glog"
	"github.com/spf13/cobra"

	osclient "github.com/projectatomic/appinfra-next/pkg/client"
	"github.com/projectatomic/appinfra-next/pkg/cmd/util/clientcmd"
	"github.com/projectatomic/appinfra-next/pkg/oauth/osintypes"
)

func NewCmdValidateToken(f *clientcmd.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate-token",
		Short: "validate an access token",
		Long:  `validate an access token`,
		Run: func(cmd *cobra.Command, args []string) {
			tokenValue := getFlagString(cmd, "token")

			clientCfg, err := f.OpenShiftClientConfig.ClientConfig()
			if err != nil {
				fmt.Errorf("%v\n", err)
			}

			validateToken(tokenValue, clientCfg)
		},
	}
	cmd.Flags().String("token", "", "Token value")
	return cmd
}

func validateToken(token string, clientConfig *kclient.Config) {
	if len(token) == 0 {
		fmt.Println("You must provide a token to validate")
		return
	}
	fmt.Printf("Using token: %v\n", token)

	clientConfig.BearerToken = token

	osClient, err := osclient.New(clientConfig)
	if err != nil {
		fmt.Printf("Error building osClient: %v\n", err)
		return
	}

	jsonResponse, _, err := getTokenInfo(token, osClient)
	if err != nil {
		fmt.Printf("%v\n", err)
		fmt.Println("Try visiting " + getRequestTokenURL(clientConfig) + " for a new token.")
		return
	}
	fmt.Printf("%v\n", string(jsonResponse))

	whoami, err := osClient.Users().Get("~")
	if err != nil {
		fmt.Printf("Error making whoami request: %v\n", err)
		return
	}
	whoamiJSON, err := json.Marshal(whoami)
	if err != nil {
		fmt.Printf("Error interpretting whoami response: %v\n", err)
		return
	}
	fmt.Printf("%v\n", string(whoamiJSON))
}

func getTokenInfo(token string, osClient *osclient.Client) (string, *osintypes.InfoResponseData, error) {
	osResult := osClient.Get().AbsPath("oauth", "info").Param("code", token).Do()
	if osResult.Error() != nil {
		return "", nil, fmt.Errorf("Error making info request: %v", osResult.Error())
	}
	body, err := osResult.Raw()
	if err != nil {
		return "", nil, fmt.Errorf("Error reading info response: %v\n", err)
	}
	glog.V(1).Infof("Raw JSON: %v\n", string(body))

	var accessData osintypes.InfoResponseData
	err = json.Unmarshal(body, &accessData)
	if err != nil {
		return "", nil, fmt.Errorf("Error while unmarshalling info response: %v %v", err, string(body))
	}
	if accessData.Error == "invalid_request" {
		return "", nil, fmt.Errorf("\"%v\" is not a valid token.\n", token)
	}
	if len(accessData.ErrorDescription) != 0 {
		return "", nil, fmt.Errorf("%v\n", accessData.ErrorDescription)
	}

	return string(body), &accessData, nil

}

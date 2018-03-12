package cmd

import (
	"github.com/spf13/cobra"
	"github.com/aws/aws-sdk-go/service/ssm"
	"encoding/json"
	"fmt"
	"strings"
	"regexp"
)

var outFormat *string
var vars *[]string
var ssmClient *ssm.SSM

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Loads provided set of parameter store params",
	Run: func(cmd *cobra.Command, args []string) {
		ssmClient = ssm.New(awsSession)
		loaded := MustLoadParams(ssmClient, *vars)

		switch *outFormat {
		case "environment":
			for k, v := range loaded {
				fmt.Printf("export %s=\"%x\"", toBashName(k), toBashValue(v))
			}
		case "json":
		default:
			st, err := json.MarshalIndent(loaded, "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(st))
			break
		}
	},
}

func init() {
	RootCommand.AddCommand(loadCmd)
	outFormat = loadCmd.Flags().StringP("output-format", "o", "", "'json', 'environment' (default 'environment')")
	vars = loadCmd.Flags().StringSliceP("vars", "v", nil, "list of vars to retrieve")
}

var escRegex = regexp.MustCompile("[^a-zA-Z0-9_]")
func toBashName(k string) string {
	return strings.ToUpper(string(escRegex.ReplaceAll([]byte(k), []byte("_"))))
}

func toBashValue(v string) string {
	if strings.Contains(v, "\n") {
		return fmt.Sprintf("$(cat <<'EOVAR'\n%s\nEOVAR\n)\n", v)
	}

	return fmt.Sprintf("'%s'", strings.Replace(v, "'", "'\\''", -1))
}

func MustLoadParam(cl *ssm.SSM, name string) string {
	res, err := LoadParam(cl, name)
	if err != nil {
		panic(err)
	}

	return res
}

func LoadParam(cl *ssm.SSM, name string) (string, error) {
	decrypt := true;
	res, err := cl.GetParameter(&ssm.GetParameterInput{Name: &name, WithDecryption: &decrypt})
	if err != nil {
		return "", err
	}

	return *res.Parameter.Value, nil
}

func LoadParams(cl *ssm.SSM, names []string) (map[string]string, error) {
	ps := []*string{}
	for _, n := range names {
		ps = append(ps, &n)
	}

	res, err := cl.GetParameters(&ssm.GetParametersInput{Names: ps})
	if err != nil {
		return nil, err
	}

	if len(res.InvalidParameters) > 0 {
		strs := []string{}
		for _, p := range res.InvalidParameters {
			strs = append(strs, *p)
		}

		return nil, fmt.Errorf("invalid parameters: %s", strings.Join(strs, ","))
	}

	vals := map[string]string{}
	for _, p := range res.Parameters {
		vals[*p.Name] = *p.Value
	}

	return vals, nil
}

func MustLoadParams(cl *ssm.SSM, names []string) map[string]string {
	res, err := LoadParams(cl, names)
	if err != nil {
		panic(err)
	}
	return res
}


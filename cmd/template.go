package cmd

import (
	"github.com/spf13/cobra"
	"github.com/aws/aws-sdk-go/service/ssm"
	"io/ioutil"
	"text/template"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"bytes"
	"fmt"
)

var src *string
var dest *string
var varPrefix *string
var params map[string]string

var tplCommand = &cobra.Command{
	Use:   "template",
	Short: "Templates the provided file from parameter store vars",
	Run: func(cmd *cobra.Command, args []string) {
		params = map[string]string{}
		ssmClient = ssm.New(awsSession)

		tplb, err := ioutil.ReadFile(*src)
		if err != nil {
			panic(err)
		}

		contents := string(tplb)
		tpl, err := template.New("").Funcs(template.FuncMap {
			"param": func(n string) (string, error) {
				return outVar(*varPrefix + n, false, "")
			},
			"paramD": func(n string, def string) (string, error) {
				return outVar(*varPrefix + n, true, def)
			},
			"paramFull": func(n string) (string, error) {
				return outVar(n, false, "")
			},
			"paramFullD": func(n string, def string) (string, error) {
				return outVar(n, true, def)
			},
		}).Parse(contents)
		if err != nil {
			panic(err)
		}

		var buf bytes.Buffer
		err = tpl.Execute(&buf, nil)
		if err != nil {
			panic(err)
		}

		if dest != nil && *dest != "" {
			ioutil.WriteFile(*dest, buf.Bytes(), 0644)
		} else {
			fmt.Print(buf.String())
		}
	},
}

func outVar(n string, allowDefault bool, def string) (string, error) {
	val, ok := params[n]
	if ok {
		return val, nil
	}

	val, err := LoadParam(ssmClient, n)
	if err != nil {
		if !allowDefault {
			return "", err
		}

		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == ssm.ErrCodeParameterNotFound {
				params[n] = def
				return def, nil
			}
		}

		return "", err
	}


	params[n] = val
	return val, nil

}

func init() {
	RootCommand.AddCommand(tplCommand)
	src = tplCommand.Flags().StringP("source", "s", "", "path to source template file")
	dest = tplCommand.Flags().StringP("dest", "d", "", "path to template to (optional, will only write on success). defaults to stdout")
	varPrefix = tplCommand.Flags().StringP("prefix", "p", "", "prefix name before sending to parameter store (e.g., 'some-app/some-env/')")
}

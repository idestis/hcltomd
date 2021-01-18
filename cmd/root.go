package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"github.com/hashicorp/hcl"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Version holds information regarding tool version
type Version struct {
	Date    string
	Version string
	Commit  string
}

var (
	showVersion bool
	cfgFile     string
	// File holds location for variable file
	File string
	// Output hold the output formatting.
	Output  string
	version Version
	// Bold is used to print text output as bold text
	Bold = color.New(color.Bold)
)

var rootShortHelp = fmt.Sprintf(
	"%v allows you to convert your terraform variables file into markdown table",
	Bold.Sprint("hcltomd"),
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hcltomd",
	Short: rootShortHelp,
	Args:  cobra.NoArgs,
	Run:   runFuncOrVersion,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(v Version) {
	version = v
	if err := rootCmd.Execute(); err != nil {
		// fmt.Printf("[x] ERROR - %v\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	versionHelpMessage := fmt.Sprintf("prints version information for %v and components", Bold.Sprint("gish"))
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, versionHelpMessage)
	rootCmd.Flags().StringVarP(&File, "file", "f", "", "terraform vars file path (required)")
	rootCmd.Flags().StringVarP(&Output, "output", "o", "table", "define the formatting of the output (e.g. table, list)")
	rootCmd.MarkFlagRequired("file")
}

func initConfig() {
	viper.SetEnvPrefix("hcltomd")
	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	// if err := viper.ReadInConfig(); err == nil {
	// 	fmt.Println("Using config file:", viper.ConfigFileUsed())
	// }
}

func runFuncOrVersion(cmd *cobra.Command, args []string) {
	if showVersion {
		fmt.Printf("%v version %s (%s) %s\n", Bold.Sprint("hcltomd"), version.Version, version.Commit, version.Date)
	} else {
		if _, err := os.Stat(File); os.IsNotExist(err) {
			logrus.Errorf("%v file is not exists", File)
			os.Exit(1)
		}
		data, err := ioutil.ReadFile(File)
		if err != nil {
			logrus.Errorf("File reading error %v", err)
			return
		}
		jsonData, err := hclToJSON(data)
		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetCenterSeparator("|")
		table.SetColumnSeparator("|")
		table.SetRowSeparator("-")
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetAutoWrapText(false)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetHeader(
			[]string{
				"Name",
				"Type",
				"Default",
				"Description",
			})
		var tableData []Variable
		for _, variable := range jsonData.([]map[string]interface{}) {
			// key == id, label, properties, etc
			for name, values := range variable {
				// Placeholders
				var def interface{}
				varType := ""
				varDesc := ""
				for _, params := range values.([]map[string]interface{}) {
					for key, param := range params {
						if key == "default" {
							def = param
						}
						if key == "description" {
							varDesc = param.(string)
						}
						if key == "type" {
							varType = param.(string)
						}
					}
				}
				tableData = append(tableData, Variable{
					Default:     def,
					Description: varDesc,
					Name:        name,
					Type:        varType,
				})
			}
		}
		for _, v := range tableData {
			s := []string{
				v.Name,
				v.Type,
				fmt.Sprintf("%v", v.Default),
				v.Description,
			}
			table.Append(s)
		}
		table.Render()
	}
}

// HashiCorp Configuration Language struct
type HashiCorp struct {
	// {"variable": []}
	Variable interface{} `json:"variable"`
	// {"output": []}
	Output interface{} `json:"output"`
}

// Variable is
type Variable struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Default     interface{} `json:"default"`
	Description string      `json:"description"`
}

func hclToJSON(content []byte) (interface{}, error) {
	var out HashiCorp
	// FIXME: Define why Decode doesn't work with type = string from TF12
	// Terraform 0.11 and earlier required type constraints to be given in quotes,
	// but that form is now deprecated and will be removed in a future version of
	// Terraform. To silence this warning, remove the quotes around "string".
	// type = "string"
	err := hcl.Unmarshal(content, &out)
	if err != nil {
		return nil, err
	}
	// TODO: allow to print outputs as well for documentation needs
	return out.Variable, nil
}

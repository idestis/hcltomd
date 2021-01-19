package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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

// HashiCorp Configuration Language struct
type HashiCorp struct {
	// {"variable": []}
	Variable interface{} `json:"variable"`
	// {"output": []}
	Output interface{} `json:"output"`
	// {"resource": []}
	Resource interface{} `json:"resource"`
	// {"data": []}
	Data interface{} `json:"data"`
}

// TableVariable is
type TableVariable struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Default     interface{} `json:"default"`
	Description string      `json:"description"`
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
	// TODO: Add more print table options, such as MD.
	// rootCmd.Flags().StringVarP(&Output, "output", "o", "table", "define the formatting of the output (e.g. table, list)")
	rootCmd.MarkFlagRequired("file")
}

func initConfig() {
	viper.SetEnvPrefix("hcltomd") // HCLTOMD_
	viper.AutomaticEnv()          // read in environment variables that match
}

func runFuncOrVersion(cmd *cobra.Command, args []string) {
	if showVersion {
		fmt.Printf("%v version %s (%s) %s\n", Bold.Sprint("hcltomd"), version.Version, version.Commit, version.Date)
	} else {
		if _, err := os.Stat(File); os.IsNotExist(err) {
			logrus.Errorf("%v file is not exists", File)
			os.Exit(1)
		}
		fileData, err := readFileAndFormat(File)
		if err != nil {
			logrus.Error(err)
			return
		}
		data, err := hclToInterface(fileData)
		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
		tableRender(data)
	}
}

func tableRender(data interface{}) {
	table := tablewriter.NewWriter(os.Stdout)
	// Table style configuration
	table.SetCenterSeparator("|")
	table.SetColumnSeparator("|")
	table.SetRowSeparator("-")
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetAutoWrapText(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	// Table Header definition
	table.SetHeader(
		[]string{
			"Name",
			"Type",
			"Default",
			"Description",
		})
	// tableData should be slice of slices,
	// but first we need to organize our data from file.
	var tableData []TableVariable
	for _, variable := range data.([]map[string]interface{}) {
		for name, values := range variable {
			// Placeholders for data
			var defaultValue interface{}
			typeValue := ""
			descriptionValue := ""
			for _, params := range values.([]map[string]interface{}) {
				for key, param := range params {
					if key == "default" {
						defaultValue = param
					}
					if key == "description" {
						descriptionValue = param.(string)
					}
					if key == "type" {
						typeValue = param.(string)
					}
				}
			}
			tableData = append(tableData, TableVariable{
				Default:     defaultValue,
				Description: descriptionValue,
				Name:        name,
				Type:        typeValue,
			})
		}
	}

	// Loop over structured data to create simple slice of slices.
	for _, v := range tableData {
		s := []string{
			// Order here is required to match the table header
			v.Name,
			v.Type,
			// Convert interface(whatever) to string
			fmt.Sprintf("%v", v.Default),
			v.Description,
		}
		table.Append(s)
	}
	table.Render()
}

func hclToInterface(content []byte) (interface{}, error) {
	var out HashiCorp
	// var data interface{}
	// FIXME: Define why Unmarshal doesn't work with unqoted type from TF12
	// Terraform 0.11 and earlier required type constraints to be given in quotes,
	// but that form is now deprecated and will be removed in a future version of
	// Terraform. To silence this warning, remove the quotes around "string".
	// type = "string"

	err := hcl.Unmarshal(content, &out)
	if err != nil {
		return nil, err
	}
	// TODO: Work on this side to select/deleselect what you want to parse.
	// Variables, Outputs, Data.
	return out.Variable, nil
}

func readFileAndFormat(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		logrus.Fatal(err)
	}
	defer file.Close()
	var out []byte
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "type") && strings.Contains(scanner.Text(), "=") {
			s := strings.Split(scanner.Text(), "=")
			if !strings.Contains(s[1], "\"") {
				s[1] = writeQuote(s[1])
				f := strings.Join(s, "= ")
				out = append(out, f...)
			} else {
				out = append(out, scanner.Text()...)
			}
		} else {
			out = append(out, scanner.Text()...)
		}
	}

	if err := scanner.Err(); err != nil {
		logrus.Fatal(err)
	}

	return out, nil
}

func writeQuote(s string) string {
	// TODO: Refactor this part, here is good example
	// https://github.com/hashicorp/terraform/blob/1c7c53adbbf7c33713c1cf6696215132cf7eaf90/command/fmt.go#L430
	s = strings.TrimSpace(s)
	// if strings.HasPrefix(s, "string") || strings.HasPrefix(s, "bool") || strings.HasPrefix(s, "number") {
	// 	return `"` + s + `"`
	// }
	return `"` + s + `"`
}

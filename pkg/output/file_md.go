package output

import (
	"html/template"
	"os"
	"time"

	"github.com/ffuf/ffuf/pkg/ffuf"
)

const (
	markdownTemplate = `# FFUF Report

  Command line : ` + "`{{.CommandLine}}`" + `
  Time: ` + "{{ .Time }}" + `

  {{ range .Keys }}| {{ . }} {{ end }}| Position | Status Code | Content Length | Content Words | Content Lines |
  {{ range .Keys }}| :- {{ end }}| :---- | :------- | :---------- | :------------- | :------------ | :------------ |
  {{range .Results}}{{ range $keyword, $value := .Input }}| {{ $value | printf "%s" }} {{ end }}| {{ .Position }} | {{ .StatusCode }} | {{ .ContentLength }} | {{ .ContentWords }} | {{ .ContentLines }} |
  {{end}}` // The template format is not pretty but follows the markdown guide
)

func writeMarkdown(config *ffuf.Config, res []Result) error {

	ti := time.Now()

	keywords := make([]string, 0)
	for _, inputprovider := range config.InputProviders {
		keywords = append(keywords, inputprovider.Keyword)
	}

	outMD := htmlFileOutput{
		CommandLine: config.CommandLine,
		Time:        ti.Format(time.RFC3339),
		Results:     res,
		Keys:        keywords,
	}

	f, err := os.Create(config.OutputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	templateName := "output.md"
	t := template.New(templateName).Delims("{{", "}}")
	t.Parse(markdownTemplate)
	t.Execute(f, outMD)
	return nil
}

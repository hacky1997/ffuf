package output

import (
	"html/template"
	"os"
	"time"

	"github.com/ffuf/ffuf/pkg/ffuf"
)

type htmlFileOutput struct {
	CommandLine string
	Time        string
	Keys        []string
	Results     []Result
}

const (
	htmlTemplate = `
<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1, maximum-scale=1.0"
    />
    <title>FFUF Report - </title>

    <!-- CSS  -->
    <link
      href="https://fonts.googleapis.com/icon?family=Material+Icons"
      rel="stylesheet"
    />
    <link
      rel="stylesheet"
      href="https://cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/css/materialize.min.css"
    />
  </head>

  <body>
    <nav>
      <div class="nav-wrapper">
        <a href="#" class="brand-logo">FFUF</a>
        <ul id="nav-mobile" class="right hide-on-med-and-down">
        </ul>
      </div>
    </nav>

    <main class="section no-pad-bot" id="index-banner">
      <div class="container">
        <br /><br />
        <h1 class="header center ">FFUF Report</h1>
        <div class="row center">

		<pre>{{ .CommandLine }}</pre>
		<pre>{{ .Time }}</pre>

   <table>
        <thead>
        <div style="display:none">
|result_raw|StatusCode|Input|Position|ContentLength|ContentWords|ContentLines|
        </div>
          <tr>
              <th>Status</th>
{{ range .Keys }}              <th>{{ . }}</th>
{{ end }}
              <th>Position</th>
              <th>Length</th>
              <th>Words</th>
              <th>Lines</th>
          </tr>
        </thead>

        <tbody>
			{{range $result := .Results}}
                <div style="display:none">
|result_raw|{{ $result.StatusCode }}{{ range $keyword, $value := $result.Input }}|{{ $value | printf "%s" }}{{ end }}|{{ $result.Position }}|{{ $result.ContentLength }}|{{ $result.ContentWords }}|{{ $result.ContentLines }}|
                </div>
				<tr class="result-{{ $result.StatusCode }}" style="background-color: {{$result.HTMLColor}};"><td><font color="black" class="status-code">{{ $result.StatusCode }}</font></td>{{ range $keyword, $value := $result.Input }}<td>{{ $value | printf "%s" }}</td>{{ end }}</td><td>{{ $result.Position }}</td><td>{{ $result.ContentLength }}</td><td>{{ $result.ContentWords }}</td><td>{{ $result.ContentLines }}</td></tr>
            {{end}}
        </tbody>
      </table>

        </div>
        <br /><br />
      </div>
    </main>

    <!--JavaScript at end of body for optimized loading-->
    <script src="https://cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/js/materialize.min.js"></script>
    <style>
      body {
        display: flex;
        min-height: 100vh;
        flex-direction: column;
      }

      main {
        flex: 1 0 auto;
      }
    </style>
  </body>
</html>

	`
)

// colorizeResults returns a new slice with HTMLColor attribute
func colorizeResults(results []Result) []Result {
	newResults := make([]Result, 0)

	for _, r := range results {
		result := r
		result.HTMLColor = "black"

		s := result.StatusCode

		if s >= 200 && s <= 299 {
			result.HTMLColor = "#adea9e"
		}

		if s >= 300 && s <= 399 {
			result.HTMLColor = "#bbbbe6"
		}

		if s >= 400 && s <= 499 {
			result.HTMLColor = "#d2cb7e"
		}

		if s >= 500 && s <= 599 {
			result.HTMLColor = "#de8dc1"
		}

		newResults = append(newResults, result)
	}

	return newResults
}

func writeHTML(config *ffuf.Config, results []Result) error {

	results = colorizeResults(results)

	ti := time.Now()

	keywords := make([]string, 0)
	for _, inputprovider := range config.InputProviders {
		keywords = append(keywords, inputprovider.Keyword)
	}

	outHTML := htmlFileOutput{
		CommandLine: config.CommandLine,
		Time:        ti.Format(time.RFC3339),
		Results:     results,
		Keys:        keywords,
	}

	f, err := os.Create(config.OutputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	templateName := "output.html"
	t := template.New(templateName).Delims("{{", "}}")
	t.Parse(htmlTemplate)
	t.Execute(f, outHTML)
	return nil
}

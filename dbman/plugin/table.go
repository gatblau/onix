//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package plugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// generic table used as a serializable result set for queries
type Table struct {
	Header Row   `json:"header,omitempty"`
	Rows   []Row `json:"row,omitempty"`
}

// a row in the table
type Row []string

// save the table to a file with the specified format
//   - filename: the filename with no extension
//   - format: either JSON, YAML/YML or CSV
func (table *Table) Save(format string, filename string) {
	// get the path of the current executing process
	ex, err := os.Executable()
	if err != nil {
		fmt.Printf("!!! I cannot find the path to the current process: %s\n", err)
	}
	exPath := filepath.Dir(ex)
	// create a file with the getReleaseInfo getPlan
	f, err := os.Create(fmt.Sprintf("%v/%v.%v", exPath, filename, format))
	if err != nil {
		fmt.Printf("!!! I cannot create the result file: %s\n", err)
	}
	f.WriteString(table.Sprint(format))
	f.Close()
}

// return the table as a string of the specified format
//   - format: either JSON, YAML/YML or CSV
func (table *Table) Sprint(format string) string {
	switch strings.ToLower(format) {
	case "yml":
		fallthrough
	case "yaml":
		return table.AsYAML()
	case "json":
		return table.AsJSON()
	case "csv":
		return table.AsCSV()
	default:
		fmt.Printf("!!! output format %v not supported, try YAML, JSON or CSV", format)
	}
	return ""
}

// return the table as a JSON string
func (table *Table) AsJSON() string {
	o, err := json.MarshalIndent(table, "", " ")
	if err != nil {
		fmt.Printf("!!! cannot convert output to JSON: %v", err)
	}
	return string(o)
}

// return the table as a YAML string
func (table *Table) AsYAML() string {
	o, err := yaml.Marshal(table)
	if err != nil {
		fmt.Printf("!!! cannot convert output to YAML: %v", err)
	}
	return string(o)
}

// return the table as a CSV string
func (table *Table) AsCSV() string {
	buffer := bytes.Buffer{}
	for i := 0; i < len(table.Header); i++ {
		buffer.WriteString(table.Header[i])
		if i < len(table.Header)-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("\n")
	for _, row := range table.Rows {
		for i := 0; i < len(row); i++ {
			buffer.WriteString(row[i])
			if i < len(row)-1 {
				buffer.WriteString(",")
			}
		}
		buffer.WriteString("\n")
	}
	out := buffer.String()
	return out[:len(out)-1]
}

// writes the table as an html page to the passed-in writer
// writer: the output stream where the html representation of the table will be written
// vars: the variables to merge when creating the html representation of the table
func (table *Table) AsHTML(writer io.Writer, vars *HtmlTableVars) error {
	t, err := template.New("report").Parse(htmlTableTemplate)
	if err != nil {
		return err
	}
	return t.Execute(writer, vars)
}

// print the table content with the specified format to the stdout
//   - format: either JSON, YAML/YML or CSV
func (table *Table) Print(format string) {
	fmt.Println(table.Sprint(format))
}

// an html template to render a Table as an html page
// merges the HtmlTableVars struct
const htmlTableTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{.Title}}</title>
    <style>
        body { font-family: Avenir }
        #tableWrap {
            display: grid;
            /*grid-template-columns: auto auto auto auto auto;*/
        }
		#title { padding: 10px; font-size: x-large; font-weight: bold; }
		#description { padding: 20px; font-style: italic; }
		#dbman { padding: 20px; font-style: italic; float:right; }
        div.cell { padding: 10px; }
        div.head {
            background: #0073ff;
            color: #fff;
            font-weight: bold;
        }
        div.alt {
            background: #f2f2f2;
        }
        /* responsive transform */
        @media screen and (max-width: 600px) {
            #tableWrap {
                grid-template-columns: 100%;
            }
            div.cell {
                padding: 5px;
            }
        }
		{{if .Style}}
		/* override base styles here */
		{{.Style}}
		{{end}}
    </style>
</head>
<body>
    {{if .Header}}
    {{.Header}}
    {{end}}
    <div id="title">{{.Title}}</div>
	{{if .Description}}
	<div id="description">{{.Description}}</div>
	{{end}}
    <div id="tableWrap"></div>
    <div id="dbman">Powered by <a href="http://onix.gatblau.org" target="_blank">Onix DbMan</a></div>
    {{if .Footer}}
    {{.Footer}}
    {{end}}
    <script language="JavaScript">
    let source;
    const queryURI = '{{.QueryURI}}';
    fetch(queryURI, {
        headers: {'Accept': 'application/json'},
    }).then(response => response.json())
      .then(source => {
        // write the table header
        const theWrap = document.getElementById("tableWrap");
        // dynamically set the number of columns
        let cols = ""
        for (let c = 0; c < source.header.length; c++) {
            cols += "auto"
            if (c < source.header.length - 1) {
                cols += " "
            }
        }
        theWrap.style.gridTemplateColumns = cols
        let theCell = null;
        for (let ix = 0; ix < source.header.length; ix++) {
            theCell = document.createElement("div");
            theCell.innerHTML = source.header[ix];
            theCell.classList.add("cell");
            theCell.classList.add("head");
            theWrap.appendChild(theCell);
        }
		// write the table rows
        let theRow = null;
        let altRow = false;
        for (let rowIx = 0; rowIx < source.row.length; rowIx++) {
            for (let cellIx = 0; cellIx < source.header.length; cellIx++) {
                theCell = document.createElement("div");
                theCell.innerHTML = source.row[rowIx][cellIx];
                theCell.classList.add("cell");
                if (altRow) {
                    theCell.classList.add("alt");
                }
                theWrap.appendChild(theCell);
            }
            altRow = !altRow;
        }
    })
</script>
</body>
</html>`

// provides merge data for htmlTableTemplate
type HtmlTableVars struct {
	// the table title
	Title string
	// the table description
	Description string
	// the URI of the query that retrieves the json table
	QueryURI string
	// the content of the CSS stylesheet to embed
	Style *string
	// the content of the stylesheet to embed
	Header *string
	// the content of the footer to embed
	Footer *string
}

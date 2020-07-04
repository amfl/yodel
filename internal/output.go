package yodel

import (
	"gopkg.in/yaml.v2"

	mapset "github.com/deckarep/golang-set"
	"github.com/jedib0t/go-pretty/v6/table"
)

// OutputYaml outputs a set of groups as yaml :)
func OutputYaml(groups mapset.Set, annotate func(string) string) (string, error) {
	// Convert to slice
	slice := groups.ToSlice()

	// Convert to yaml
	d, err := yaml.Marshal(&slice)
	if err != nil {
		return "", err
	}

	return string(d), nil
}

// OutputTable outputs a set of groups as a table with go-pretty
func OutputTable(groups mapset.Set) (string, error) {
	// Convert to slice
	slice := groups.ToSlice()

	tw := table.NewWriter()

	rows := make([]table.Row, len(slice))
	for i, r := range slice {
		rows[i] = table.Row{r}
	}
	// Append a header
	tw.AppendRows(rows)
	return tw.Render(), nil
}

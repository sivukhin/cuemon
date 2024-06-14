package lib

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateVariable(t *testing.T) {
	var grafana JsonRaw[Grafana]
	require.Nil(t, json.Unmarshal(readFile(path.Join("files", "example_1.json")), &grafana))
	variable := grafana.Value.Value.Templating.List[3]
	panel, err := MonitoringContext{SchemaFiles: []string{GrafanaV10Cue}}.CreateVariableSrc(variable)
	require.Nil(t, err)
	t.Log(FormatNode(File(panel)))
}

func TestCreatePanel(t *testing.T) {
	var grafana JsonRaw[Grafana]
	require.Nil(t, json.Unmarshal(readFile(path.Join("files", "example_1.json")), &grafana))
	context := MonitoringContext{SchemaFiles: []string{GrafanaV10Cue}}
	dashboard, err := context.CreateDashboardSrc(grafana.Value)
	require.Nil(t, err)
	t.Log(FormatNode(File(dashboard.Meta)))
}

func TestCreateRow(t *testing.T) {
	//var grafana JsonRaw[Grafana]
	//require.Nil(t, json.Unmarshal(readFile(path.Join("files", "example_1.json")), &grafana))
	//rows := ExtractRows(grafana.Value)
	//row, err := MonitoringContext{}.CreateRowSrc(rows[0])
	//require.Nil(t, err)
	//result, err := FormatNode(File(row.Src))
	//require.Nil(t, err)
	//t.Log(result)
}

func readFile(path string) []byte {
	result, err := os.ReadFile(path)
	if err != nil {
		Logger.Fatalf("failed to read file %v: %v", path, err)
	}
	return result
}

package stanza_input_operator

import (
	"context"
	"encoding/json"
	"net/http"
)

type NodeField struct {
	FieldName   string `json:"field_name"`
	Type        string `json:"type"`
	Color       string `json:"color,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

type EdgeField struct {
	FieldName string `json:"field_name"`
	Type      string `json:"type"`
}

type GraphFields struct {
	NodesFields []NodeField `json:"nodes_fields"`
	EdgesFields []EdgeField `json:"edges_fields"`
}
type Node struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	SubTitle string `json:"subTitle"`
	MainStat string `json:"mainStat"`
	// DetailRole string `json:"detail__role,omitempty"`
}

type Edge struct {
	ID       string `json:"id"`
	Source   string `json:"source"`
	Target   string `json:"target"`
	MainStat string `json:"mainStat"`
}
type GraphData struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

var nodes []Node
var edges []Edge

func fetchGraphFields(w http.ResponseWriter, r *http.Request) {
	nodesFields := []NodeField{
		{"id", "string", "", ""},
		{"title", "string", "", ""},
		{"subTitle", "string", "", ""},
		{"mainStat", "string", "", ""},
		// {"secondaryStat", "number", "", ""},
		// {"arc__failed", "number", "red", "Failed"},
		// {"arc__passed", "number", "green", "Passed"},
		// {"detail__role", "string", "", "Role"},
	}
	edgesFields := []EdgeField{
		{"id", "string"},
		{"source", "string"},
		{"target", "string"},
		{"mainStat", "string"},
	}
	graphFields := GraphFields{NodesFields: nodesFields, EdgesFields: edgesFields}
	json.NewEncoder(w).Encode(graphFields)
}

func checkHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API is working well!"))
}

func fetchGraphData(w http.ResponseWriter, r *http.Request) {

	graphData := GraphData{Nodes: nodes, Edges: edges}
	json.NewEncoder(w).Encode(graphData)

}

func (operator *Input) nodeGraph(ctx context.Context) {
	logFilter := operator.LogFilter
	if logFilter != "all" && logFilter != "system" {
		return
	}
	nodes = []Node{}
	edges = []Edge{}
	operator.wg.Add(1)
	go func() {
		defer operator.wg.Done()

		for {
			err, log := operator.kubearmorClient.recvLogs(operator)
			if err != nil {
				operator.Logger().Warnf("%s", err.Error())
				return
			}

			if log.TTY != "" {
				node := Node{
					ID:       string(log.HostPID),
					Title:    log.ProcessName,
					SubTitle: log.PodName,
					MainStat: log.Resource,
				}
				nodes = append(nodes, node)
				if log.PPID != 0 {
					edge := Edge{
						ID:       string(log.HostPID),
						Source:   string(log.HostPPID),
						Target:   string(log.HostPID),
						MainStat: log.Resource,
					}
					edges = append(edges, edge)
				}

			}

			select {
			case <-ctx.Done():
				return
			default:
				continue
			}
		}
	}()
}

package grafananodegraphexporter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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

func (nge *nodeGraphExporter) nodeGraph(ctx context.Context) {
	defer nge.wg.Done()
	logFilter := nge.config.LogFilter
	if logFilter != "all" && logFilter != "system" {
		return
	}
	nodes = []Node{}
	edges = []Edge{}
	nge.wg.Add(1)
	go func() {
		defer nge.wg.Done()

		for {
			err, log := nge.kubearmorClient.recvLogs()
			if err != nil {
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

func (nge *nodeGraphExporter) startServer(ctx context.Context) {
	defer nge.wg.Done()

	r := mux.NewRouter()

	r.HandleFunc("/api/graph/fields", fetchGraphFields).Methods("GET")
	r.HandleFunc("/api/graph/data", fetchGraphData).Methods("GET")
	r.HandleFunc("/api/health", checkHealth).Methods("GET")

	srv := &http.Server{
		Addr:    ":5000",
		Handler: r,
	}

	// Run the server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Errorf("server listen: %s\n", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()
}

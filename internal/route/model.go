package route

import "errors"

// Represents a city or waypoint in the graph e.g. "KTM", "PKR", "BRT".
type NodeID string

// Edge connects two nodes with a cost (distance, time, fuel).
// Directed — from To has a direction, the graph can be one-way.
type Edge struct {
	To   NodeID
	Cost int
}

// Graph is an adjacency list — the standard representation for Dijkstra.
// map key = a node, value = all edges leaving that node.
// e.g. {"KTM": [{To:"PKR", Cost:200}, {To:"BRT", Cost:150}]}
type Graph map[NodeID][]Edge

// Route is the result of a shortest path query.
type Route struct {
	From      NodeID   `json:"from"`
	To        NodeID   `json:"to"`
	Path      []NodeID `json:"path"`       // ordered list of nodes to travel
	TotalCost int      `json:"total_cost"` // sum of edge costs along the path
}

// Sentinel errors for the route domain.
var (
	ErrNodeNotFound = errors.New("node not found in graph")
	ErrNoPathExists = errors.New("no path exists between nodes")
)

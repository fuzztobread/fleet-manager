package route

import (
	"errors"
	"fmt"
)

// Service holds the business logic for the route domain.
type Service struct {
	graph Graph
}

// NewService builds the service with a pre-seeded graph.
func NewService() *Service {
	// Adjacency list — directed weighted graph.
	// Costs represent distance in km between Nepali cities.
	g := Graph{
		"KTM": {{To: "PKR", Cost: 200}, {To: "BRT", Cost: 150}},
		"PKR": {{To: "KTM", Cost: 200}, {To: "BRT", Cost: 100}, {To: "JMP", Cost: 180}},
		"BRT": {{To: "KTM", Cost: 150}, {To: "PKR", Cost: 100}, {To: "JMP", Cost: 220}},
		"JMP": {{To: "PKR", Cost: 180}, {To: "BRT", Cost: 220}},
	}

	return &Service{graph: g}
}

// FindRoute finds the shortest path between two nodes.
func (s *Service) FindRoute(from, to string) (*Route, error) {
	if from == "" || to == "" {
		return nil, errors.New("from and to are required")
	}

	if from == to {
		return nil, errors.New("from and to must be different nodes")
	}

	route, err := s.graph.Shortest(NodeID(from), NodeID(to))
	if err != nil {
		// surface domain errors directly — handler will map them to HTTP status
		return nil, fmt.Errorf("find route: %w", err)
	}

	return route, nil
}

// Nodes returns all known nodes in the graph.
func (s *Service) Nodes() []NodeID {
	nodes := make([]NodeID, 0, len(s.graph))
	for node := range s.graph {
		nodes = append(nodes, node)
	}
	return nodes
}

package route

import "container/heap"

// We need the node ID and the cost to reach it from the source.
type item struct {
	node NodeID
	cost int
}

// pq is a slice of items that implements heap.Interface.
type pq []*item

func (q pq) Len() int           { return len(q) }
func (q pq) Less(i, j int) bool { return q[i].cost < q[j].cost } // min-heap: lowest cost first
func (q pq) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }

func (q *pq) Push(x any) {
	*q = append(*q, x.(*item))
}

func (q *pq) Pop() any {
	old := *q
	n := len(old)
	x := old[n-1]  // take the last element (heap.Pop swaps min to end before calling this)
	*q = old[:n-1] // shrink the slice
	return x
}

// Shortest finds the lowest-cost path from src to dst in the graph.
// Returns a Route with the full path and total cost.
// Returns ErrNodeNotFound if either node isn't in the graph.
// Returns ErrNoPathExists if dst is unreachable from src.
func (g Graph) Shortest(src, dst NodeID) (*Route, error) {
	// guard: both nodes must exist in the graph
	if _, ok := g[src]; !ok {
		return nil, ErrNodeNotFound
	}
	if _, ok := g[dst]; !ok {
		return nil, ErrNodeNotFound
	}

	// dist tracks the best known cost to reach each node.
	// Start with "infinity" (a very large number) for all nodes —
	// we haven't found a path to any of them yet.
	const inf = int(^uint(0) >> 1) // max int, platform-independent
	dist := make(map[NodeID]int)
	for node := range g {
		dist[node] = inf
	}
	dist[src] = 0 // cost to reach the source from itself is zero

	// prev records which node we came from to reach each node.
	prev := make(map[NodeID]NodeID)

	// initialise the priority queue with just the source node
	q := &pq{{node: src, cost: 0}}
	heap.Init(q) // satisfies heap invariant (no-op for single element, but correct)

	for q.Len() > 0 {
		curr := heap.Pop(q).(*item)

		// if we've reached the destination, stop early
		if curr.node == dst {
			break
		}

		// skip if we already found a cheaper path to this node
		// (stale entries can remain in the heap from earlier iterations)
		if curr.cost > dist[curr.node] {
			continue
		}

		// relax each edge leaving the current node
		for _, edge := range g[curr.node] {
			newCost := dist[curr.node] + edge.Cost

			// only update if this path is cheaper than what we know
			if newCost < dist[edge.To] {
				dist[edge.To] = newCost
				prev[edge.To] = curr.node // remember we came from curr.node
				heap.Push(q, &item{node: edge.To, cost: newCost})
			}
		}
	}

	// if dist[dst] is still infinity, dst was never reached
	if dist[dst] == inf {
		return nil, ErrNoPathExists
	}

	// reconstruct path by walking prev map backwards from dst to src
	path := buildPath(prev, src, dst)

	return &Route{
		From:      src,
		To:        dst,
		Path:      path,
		TotalCost: dist[dst],
	}, nil
}

// buildPath walks the prev map from dst back to src,
// then reverses the result so it reads src → ... → dst.
func buildPath(prev map[NodeID]NodeID, src, dst NodeID) []NodeID {
	path := []NodeID{}

	// start at destination and follow prev pointers back to source
	for node := dst; node != src; node = prev[node] {
		path = append(path, node)
	}
	path = append(path, src) // add the source node itself

	// path is currently [dst, ..., src] — reverse it
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}

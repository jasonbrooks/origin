package graph

import (
	"fmt"

	"github.com/gonum/graph"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
)

var (
	UnknownNodeKind = "UnknownNode"
)

var (
	UnknownEdgeKind = "UnknownEdge"
	// ReferencedByEdgeKind is the kind to use if you're building reverse links that don't have a specific edge in the other direction
	// other uses are discouraged.  You should create a kind for your edge
	ReferencedByEdgeKind = "ReferencedBy"
	// ContainsEdgeKind is the kind to use if one node's contents logically contain another node's contents.  A given node can only have
	// a single inbound Contais edge.  The code does not prevent contains cycles, but that's insane, don't do that.
	ContainsEdgeKind = "Contains"
)

func GetUniqueRuntimeObjectNodeName(nodeKind string, obj runtime.Object) UniqueName {
	meta, err := kapi.ObjectMetaFor(obj)
	if err != nil {
		panic(err)
	}

	return UniqueName(fmt.Sprintf("%s|%s/%s", nodeKind, meta.Namespace, meta.Name))
}

// GetTopLevelContainerNode traverses the reverse ContainsEdgeKind edges until it finds a node
// that does not have an inbound ContainsEdgeKind edge.  This could be the node itself
func GetTopLevelContainerNode(g Graph, containedNode graph.Node) graph.Node {
	// my kingdom for a LinkedHashSet
	visited := map[int]bool{}
	prevContainingNode := containedNode

	for {
		visited[prevContainingNode.ID()] = true
		currContainingNode := GetContainingNode(g, prevContainingNode)

		if currContainingNode == nil {
			return prevContainingNode
		}
		if _, alreadyVisited := visited[currContainingNode.ID()]; alreadyVisited {
			panic(fmt.Sprintf("contains cycle in %v", visited))
		}

		prevContainingNode = currContainingNode
	}

	// can't happen
	panic(fmt.Sprintf("math failed %v", visited))
	return nil
}

// GetContainingNode returns the direct predecessor that is linked to the node by a ContainsEdgeKind.  It returns
// nil if no container is found.
func GetContainingNode(g Graph, containedNode graph.Node) graph.Node {
	for _, node := range g.Predecessors(containedNode) {
		edge := g.EdgeBetween(node, containedNode)

		if g.EdgeKinds(edge).Has(ContainsEdgeKind) {
			return node
		}
	}

	return nil
}

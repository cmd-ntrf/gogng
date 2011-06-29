package main

import "flag"
import "fmt"
import "json"
import "math"
import "os"
import "rand"

type Node struct {
	error     float64
	point     []float64
	neighbors map[*Node]*Edge
}

func NewNode(point []float64, error float64) (node *Node) {
	node = new(Node)
	node.point = point
	node.error = error
	node.neighbors = make(map[*Node]*Edge)
	return
}

type Edge struct {
	age     uint
	vertex1 *Node
	vertex2 *Node
}

type Graph struct {
	edges map[*Edge]bool
	nodes map[*Node]bool
}

func NewGraph() (graph *Graph) {
	graph = new(Graph)
	graph.edges = make(map[*Edge]bool)
	graph.nodes = make(map[*Node]bool)
	return
}

func (this *Graph) AddEdge(vertex1, vertex2 *Node) (outEdge *Edge) {
	// Verify if there's an edge between the two vertices
	if edge, ok := vertex1.neighbors[vertex2]; ok {
		edge.age = 0
		outEdge = edge
		return
	}

	// Add the nodes that were not present in the graph
	if _, ok := this.nodes[vertex1]; !ok {
		this.nodes[vertex1] = true
	}
	if _, ok := this.nodes[vertex2]; !ok {
		this.nodes[vertex2] = true
	}

	// Add the new edge
	outEdge = &Edge{vertex1: vertex1, vertex2: vertex2, age: 0}
	vertex1.neighbors[vertex2] = outEdge
	vertex2.neighbors[vertex1] = outEdge
	this.edges[outEdge] = true
	return
}

func (this *Graph) RemoveEdge(edge *Edge) {

	if _, ok := this.edges[edge]; !ok {
		return
	}

	vertex1 := edge.vertex1
	vertex2 := edge.vertex2

	vertex1.neighbors[vertex2] = nil, false
	vertex2.neighbors[vertex1] = nil, false

	if len(vertex1.neighbors) == 0 {
		this.nodes[vertex1] = false, false
	}
	if len(vertex2.neighbors) == 0 {
		this.nodes[vertex2] = false, false
	}
	this.edges[edge] = false, false
	return
}

func (this *Graph) MarshalJSON() ([]byte, os.Error) {
	var output string
	output += fmt.Sprintln("{")
	output += fmt.Sprintln("\t\"nodes\":")
	output += fmt.Sprintln("\t{")
	counter := 0
	for node := range this.nodes {
		counter++
		output += fmt.Sprintf("\t\t\"%p\" : [", node)
		for idx, value := range node.point {
			output += fmt.Sprintf("%v", value)
			if idx < len(node.point)-1 {
				output += fmt.Sprintf(", ")
			}
		}
		output += fmt.Sprintf("]")
		if counter != len(this.nodes) {
			output += fmt.Sprintln(",")
		}
	}
	output += fmt.Sprintln("\n\t},")
	output += fmt.Sprintln("\t\"edges\":")
	counter = 0
	for edge := range this.edges {
		counter++
		if counter != 1 {
			output += fmt.Sprintf("\t\t[\"%p\", \"%p\"]", edge.vertex1, edge.vertex2)
		} else {
			output += fmt.Sprintf("\t\t[[\"%p\", \"%p\"]", edge.vertex1, edge.vertex2)
		}
		if counter != len(this.edges) {
			output += fmt.Sprintln(",")
		}
	}
	output += fmt.Sprintln("]")
	output += fmt.Sprintln("}")

	return []byte(output), nil
}

func Signal() []float64 {
	x := rand.Float64()*20 - 10
	var y float64
	if x > 0 {
		y = math.Sin(x) + 2.5
	} else {
		y = math.Sin(x) - 2.5
	}
	return []float64{x, y}
}

func main() {
	var lTmax = flag.Uint("tmax", 1000, "Maximum number of iterations.")
	var lTau = flag.Uint("tau", 100, "Number of iterations between two insertion.")
	var lEthag = flag.Float64("ethag", 0.2, "Winner learning rate.")
	var lEthav = flag.Float64("ethav", 0.006, "Winner's neighbors learning rate.")
	var lAmax = flag.Uint("amax", 50, "Maximum edge's age.")
	var lAlpha = flag.Float64("alpha", 0.5, "Winner forgetting rate.")
	var lDelta = flag.Float64("delta", 0.995, "Forgetting rate.")
	var lFilename = flag.String("file", "", "Resulting graph output file.")
	flag.Parse()

	lGNG := NewGraph()
	node1 := NewNode([]float64{rand.Float64(), rand.Float64()}, 0)
	node2 := NewNode([]float64{rand.Float64(), rand.Float64()}, 0)
	lGNG.AddEdge(node1, node2)

	for t := uint(1); t <= *lTmax; t++ {
		signal := Signal()

		var g1, g2 *Node
		min1, min2 := math.MaxFloat64, math.MaxFloat64

		// Find the 2 nodes closest to the signal
		for node := range lGNG.nodes {
			var error float64
			for idx, value := range signal {
				error += (node.point[idx] - value) * (node.point[idx] - value)
			}
			switch {
			case error < min1:
				g2, min2 = g1, min1
				g1, min1 = node, error
			case error < min2:
				g2, min2 = node, error
			}
		}

		// Increment adjacent edges adge
		for _, edge := range g1.neighbors {
			edge.age++
		}

		// Increment winner error
		g1.error += math.Sqrt(min1)

		// Move the adjacent nodes towards the signal
		for idx, value := range signal {
			g1.point[idx] += (*lEthag) * (value - g1.point[idx])
		}
		for node := range g1.neighbors {
			for idx, value := range signal {
				node.point[idx] += (*lEthav) * (value - node.point[idx])
			}
		}

		// Add the edge between the two nodes, if it exists, the age is just refreshed
		lGNG.AddEdge(g1, g2)

		// Remove the edges that are too old
		for edge := range lGNG.edges {
			if edge.age > *lAmax {
				lGNG.RemoveEdge(edge)
			}
		}

		// Add a node if it is the right time
		if t%*lTau == 0 {
			var q, r, x *Node
			max := -math.MaxFloat64
			for node := range lGNG.nodes {
				if node.error > max {
					max = node.error
					q = node
				}
			}
			max = -math.MaxFloat64
			for node := range q.neighbors {
				if node.error > max {
					max = node.error
					r = node
				}
			}

			lGNG.RemoveEdge(q.neighbors[r])

			point := make([]float64, len(signal))
			for idx := range signal {
				point[idx] = (q.point[idx] + r.point[idx]) / 2.0
			}
			q.error *= *lAlpha
			r.error *= *lAlpha
			x = NewNode(point, q.error)
			lGNG.AddEdge(q, x)
			lGNG.AddEdge(r, x)
		}

		// Reduce node error
		for node := range lGNG.nodes {
			node.error *= *lDelta
		}
	}

	if *lFilename != "" {
		// Outputs the resulting nodes and edges in a JSON dictionary for plotting
		file, _ := os.Open(*lFilename, os.O_WRONLY|os.O_CREAT|os.O_TRUNC, 0655)
		defer file.Close()
		encoder := json.NewEncoder(file)
		encoder.Encode(lGNG)
	}
}

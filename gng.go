package main

import "encoding/csv"
import "encoding/json"
import "flag"
import "fmt"
import "io"
import "log"
import "math"
import "math/rand"
import "os"
import "strconv"

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

func NewRandomNode(ndim int) (*Node){
	point := make([]float64, ndim)
	for i := 0; i < ndim; i++ {
		point[i] = rand.Float64()
	}
	return NewNode(point, 0)
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

	// vertex1.neighbors[vertex2] = nil, false
	delete(vertex1.neighbors, vertex2)
	delete(vertex2.neighbors, vertex1)

	if len(vertex1.neighbors) == 0 {
		delete(this.nodes, vertex1)
	}
	if len(vertex2.neighbors) == 0 {
		delete(this.nodes, vertex2)
	}
	delete(this.edges, edge)
	return
}

func (this *Graph) MarshalJSON() ([]byte, error) {
	graph := make(map[string] interface{})
	nodes := make(map[string] []float64)
	errors := make(map[string] float64)
	for node := range this.nodes {
		nodes[fmt.Sprintf("%p", node)] = node.point
		errors[fmt.Sprintf("%p", node)] = node.error
	}
	edges := make(map[string] [2]string)
	ages := make(map[string] uint)
	for edge := range this.edges {
		var edge_name string = fmt.Sprintf("%p", edge)
		edges[edge_name] = [2]string{fmt.Sprintf("%p", edge.vertex1), fmt.Sprintf("%p", edge.vertex2)}
		ages[edge_name] = edge.age 
	}
	graph["nodes"] = nodes
	graph["errors"] = errors
	graph["edges"] = edges
	graph["ages"] = ages
	return json.Marshal(graph)
}

func (this *Graph) UnmarshalJSON(input []byte) (error) {
	graph_map := make(map[string] interface{})
	json.Unmarshal(input, &graph_map)
	nodes_map := graph_map["nodes"].(map[string] interface{})
	errors_map := graph_map["errors"].(map[string] interface{})
	edges_map := graph_map["edges"].(map[string] interface{})
	ages_map := graph_map["ages"].(map[string] interface{})
	nodes := make(map[string] *Node)

	for node, point_itf := range nodes_map {
		point := make([]float64, len(point_itf.([]interface{})))
		for i, value := range point_itf.([]interface{}) {
			point[i] = value.(float64)
		}
		nodes[node] = NewNode(point, 0)
		nodes[node].error = errors_map[node].(float64)
	}
	for edge, vertex := range edges_map {
		node1 := nodes[vertex.([]interface{})[0].(string)]
		node2 := nodes[vertex.([]interface{})[1].(string)]
		new_edge := this.AddEdge(node1, node2)
		new_edge.age = uint(ages_map[edge].(float64))
	}
	return nil
}

func Signal(reader *csv.Reader) ([]float64, error) {
	fields, err := reader.Read()
	if err != nil {
		return nil, err
	}
	ndim := len(fields)
	point := make([]float64, ndim)
	for i := 0; i < ndim; i++ {
		point[i], err = strconv.ParseFloat(fields[i], 64)
		if err != nil {
			return nil, err
		}
	}
	return point, nil
}

func main() {
	var lTau = flag.Uint("tau", 100, "Number of iterations between two insertion.")
	var lOperiod = flag.Uint("operiod", 0, "Period of the topology output.")
	var lEthag = flag.Float64("ethag", 0.2, "Winner learning rate.")
	var lEthav = flag.Float64("ethav", 0.006, "Winner's neighbors learning rate.")
	var lAmax = flag.Uint("amax", 50, "Maximum edge's age.")
	var lAlpha = flag.Float64("alpha", 0.5, "Winner forgetting rate.")
	var lDelta = flag.Float64("delta", 0.995, "Forgetting rate.")
	var lOutput = flag.String("output", "", "Final topology output file.")
	var lInput = flag.String("input", "", "Initial topology input file.")
	var lData = flag.String("data", "", "CSV dataset filename.")
	flag.Parse()

	var file = os.Stdin
	if *lData != "" {
		var err error
		file, err = os.Open(*lData)
		defer file.Close()
		if err != nil {
			log.Fatalf("Can't open dataset file; err=%s\n", err.Error())
		}
	}
	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	signal, err := Signal(reader)
	if err != nil {
		log.Fatalf("Error while reading dataset, err=%s\n", err.Error())
	}

	lGNG := NewGraph()
	if *lInput != "" {
		var err error
		file, err = os.Open(*lInput)
		defer file.Close()
		if err != nil {
			log.Fatalf("Can't open input topology file; err=%s\n", err.Error())
		}
		decoder := json.NewDecoder(file)
		decoder.Decode(lGNG)
	} else {
		ndim := len(signal)
		node1 := NewRandomNode(ndim)
		node2 := NewRandomNode(ndim)
		lGNG.AddEdge(node1, node2)
	}

	t := uint(1)
	for {
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

		if *lOperiod != 0 && (t % *lOperiod) == 0 {
			encoder := json.NewEncoder(os.Stdout)
			encoder.Encode(lGNG)
		}

		// Retrieve next signal
		t++
		signal, err = Signal(reader)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Error while reading dataset, err=%s\n", err.Error())
		}
	}

	// Outputs the resulting nodes and edges in a JSON dictionary for plotting
	file = os.Stdout
	if *lOutput != "" {
		var err error
		file, err = os.OpenFile(*lOutput, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		defer file.Close()
		if err != nil {
			log.Fatalf("Can't open output topology file; err=%s\n", err.Error())
		}
	}
	encoder := json.NewEncoder(file)
	encoder.Encode(lGNG)
}

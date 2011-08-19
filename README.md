# gogng

This is a command-line implementation in Go of the classic Growing Neural Gas Network algorithm described in _A Growing Neural Gas Network Learns Topologies_ by Bernd Fritzke.

# Usage

# Plotting topology

If the signal is in 2-D, the resulting topology can be drawn with the provided script `plot_graph.py`.

## Usage 

	plot_graph.py [topology json file]

If the json file is not provided, the script tries to read the topology from the standard input. It is therefore possible to pipe the output of gogng directly in `plot_graph.py`.

# Examples

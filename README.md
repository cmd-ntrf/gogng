# gogng

This is a command-line implementation in Go of the classic Growing Neural Gas Network algorithm described in __A Growing Neural Gas Network Learns Topologies__ by Bernd Fritzke.

# Usage

# Plotting topology

If the signal is in 2-D, the resulting topology can be drawn with the provided script __plot_graph.py__.

## Usage 

	plot_graph.py [topology json file]

If the json file is not provided, the script tries to read the topology from the standard input. It is therefore possible to pipe the output of gogng directly in __plot_graph.py__

# Examples

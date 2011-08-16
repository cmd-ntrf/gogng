#!/usr/bin/env python2.7
import json
import networkx as nx
import matplotlib.pyplot as plt
import sys


try:
    graph = json.load(open(sys.argv[1], 'r'))
except IndexError:
    graph = json.load(sys.stdin)
except:
    print "Unable to open file."
    exit()

G = nx.Graph()
G.add_edges_from(graph["edges"])
nx.draw_networkx_edges(G,graph["nodes"],alpha=0.5,width=2)
nx.draw_networkx_nodes(G,graph["nodes"],node_size=4,node_color='r')

print "Edges: ", len(graph["edges"])
print "Nodes: ", len(graph["nodes"])

plt.show()

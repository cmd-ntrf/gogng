#!/usr/bin/env python2.7
import json
import networkx as nx
import matplotlib.pyplot as plt
import sys

try:
    fp = open(sys.argv[1], 'r')
except IndexError:
    fp = sys.stdin
except:
    print "Unable to open file."
    exit()

fig = plt.figure(figsize=(10,10))
ax = plt.subplot(111)
plt.ion()
plt.show()

G = nx.Graph()
def update_graph(fp):
    try:
        graph = json.loads(fp.readline())
    except ValueError:
        return False
    G.clear()
    ax.clear()
    G.add_edges_from(graph["edges"].itervalues())
    nx.draw_networkx_edges(G,graph["nodes"],alpha=0.5,width=2, animated=True)
    nx.draw_networkx_nodes(G,graph["nodes"],node_size=4,node_color='r',
                           animated=True)
    fig.canvas.draw()
    return True

while update_graph(fp): pass

plt.ioff()
plt.show()

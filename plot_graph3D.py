#!/usr/bin/env python2.7
import json
import networkx as nx
import numpy as np
from mayavi import mlab
import sys


try:
    graph = json.load(open(sys.argv[1], 'r'))
except IndexError:
    graph = json.load(sys.stdin)
except:
    print "Unable to open file."
    exit()

G = nx.Graph()
G.add_edges_from(graph["edges"].itervalues())

H = nx.convert_node_labels_to_integers(G)
# 3D spring layout
pos = nx.spring_layout(H, dim=3)

xyz=np.array([pos[v] for v in sorted(H)])
# scalar solors
scalars=np.array(H.nodes())+5

mlab.figure(1, bgcolor=(0,0,0))
mlab.clf()

pts = mlab.points3d(xyz[:,0], xyz[:, 1], xyz[:, 2],
                    scalars,
                    scale_factor=0.005,
                    scale_mode='none',
                    colormap='Blues',
                    resolution=20)

pts.mlab_source.dataset.lines = np.array(H.edges())
tube = mlab.pipeline.tube(pts, tube_radius=0.001)
mlab.pipeline.surface(tube, color=(0.8, 0.8, 0.8))

mlab.show()
#mlab.savefig('graph3d.png')

print "Edges: ", len(graph["edges"])
print "Nodes: ", len(graph["nodes"])


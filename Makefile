GC = 6g
GL = 6l

all: gng graph

gng: gng.go
	$(GC) gng.go
	$(GL) -o gng gng.6

graph: graph.go
	$(GC) graph.go
	$(GL) -o graph graph.6

clean6:
	rm -f *.6

clean:
	rm -f *.6 graph gng

.PHONY: clean clean6
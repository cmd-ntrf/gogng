GC = 6g
GL = 6l

all: gogng 

gogng: gng.go
	$(GC) gng.go
	$(GL) -o gogng gng.6

clean6:
	rm -f *.6

clean:
	rm -f *.6 gng

.PHONY: clean clean6

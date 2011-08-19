#!/usr/local/bin/python
import json
import sys
try:
    print json.dumps(json.load(open(sys.argv[1])), indent=True)
except IndexError:
    print json.dumps(json.load(sys.stdin), indent=True)
except:
    print "Unable to open file " + sys.argv[1] 


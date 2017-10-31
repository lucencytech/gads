#!/usr/local/bin/python3

class Node:
    color = None
    key = None

    def __init__(self, k):
        self.key = k
    def __str__(self):
        return "Node[key={}, color={}]".format(self.key, self.color)
    def __repr__(self):
        return self.__str__()

class Graph:
    edges = set()
    nodes = set()
    nodesByKey = {}
    graph = {}

    colors = list(range(1,11))

    # idempotent, as it only initializes empty edge set when necessary
    def addNode(self, node):
        self.nodes.add(node)
        self.nodesByKey[node.key] = node
        if self.graph.get(node.key) is None:
            self.graph[node.key] = set()

    # idempotent, since it just adds to sets
    def addEdge(self, n1, n2):
        self.graph[n1.key].add(n2)
        self.graph[n2.key].add(n1)

    def __str__(self):
        return str(self.graph)

def main():
    G = Graph()

    onHeader = True
    f = open('field-exclusions.CAMPAIGN_PERFORMANCE_REPORT.csv', 'r')
    for line in f:
        # skip header line
        if onHeader:
            onHeader = False
            continue
        fieldName, exclusionsString = line.strip().split(',')
        exclusions = exclusionsString.split(';')

        G.addNode(Node(fieldName))
        for exclusion in exclusions:
            if exclusion == '' or exclusion is None:
                continue
            G.addNode(Node(exclusion))
            G.addEdge(Node(fieldName), Node(exclusion))
    print("Graph:", G)

if __name__ == '__main__':
    main()

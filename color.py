#!/usr/local/bin/python3

from pprint import pprint
from operator import attrgetter

class Node:
    color = None
    key = None

    def __init__(self, k):
        self.key = k
    def __str__(self):
        return "{}<{}>".format(self.key, self.color)
    def __repr__(self):
        return self.__str__()
    def __hash__(self):
        return hash(self.key)
    def __eq__(self, other):
        return self.key == other.key

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

    def getDegree(self, node):
        return len(self.graph[node.key])

    def color(self):
        nodes = self.nodes
        ns = sorted(nodes, key=attrgetter('key'))
        sortedNodes = sorted(ns, key=self.getDegree, reverse=True)
        print("sortedNodes:", sortedNodes)

        color = 1

        while len(sortedNodes) != 0:
            node = sortedNodes.pop(0)
            node.color = color
            for candidateNode in sortedNodes:
                # Don't color nodes that already have a color
                if candidateNode.color is not None:
                    continue
                # Don't color any adjacent nodes with the same color
                if self.graph.get(node.key) == candidateNode.key:
                    continue
                candidateNode.color = color
                # sortedNodes = [n for n in sortedNodes if n.key != candidateNode.key]
            color += 1


def main():
    G = Graph()

    onHeader = True
    # f = open('field-exclusions.CAMPAIGN_PERFORMANCE_REPORT.csv', 'r')
    f = open('example-graph.csv', 'r')
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
    print("Graph:")
    pprint(G.graph)
    G.color()
    pprint(G.graph)
    pprint(G.nodes)

    return G

if __name__ == '__main__':
    main()

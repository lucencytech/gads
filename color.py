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
    # edges = set()
    # nodes = set()
    # nodesByKey = {}
    graph = {}

    colors = list(range(1,11))

    # idempotent, as it only initializes empty edge set when necessary
    def addNode(self, node):
        # self.nodes.add(node)
        # self.nodesByKey[node.key] = node
        if self.graph.get(node) is None:
            self.graph[node] = set()

    # idempotent, since it just adds to sets
    def addEdge(self, n1, n2):
        self.graph[n1].add(n2)
        self.graph[n2].add(n1)

    def __str__(self):
        return str(self.graph)

    def getDegree(self, node):
        return len(self.graph[node])

    def color(self):
        nodes = self.graph.keys()
        ns = sorted(nodes, key=attrgetter('key'))
        sortedNodes = sorted(ns, key=self.getDegree, reverse=True)
        print("sortedNodes:", sortedNodes)

        thisColor = 1

        while len(sortedNodes) != 0:
            node = sortedNodes.pop(0)
            node.color = thisColor
            nodesWithThisColor = [node]
            print("Colored", node, "with color", thisColor)
            for foo in sortedNodes:
                shouldColor = True # assume we should color it until we find otherwise
                candidateNode = [k for k in self.graph.keys() if k.key == foo.key][0]
                # print("node:", node)
                # print("foo:", foo)
                print("candidateNode:", candidateNode)
                # print("node's adjacent nodes:", self.graph.get(node))

                # Don't color nodes that already have a color
                if candidateNode.color is not None:
                    print("Not recoloring", candidateNode)
                    shouldColor = False
                    continue
                # Don't color any adjacent nodes with the same color
                for nodeWithThisColor in nodesWithThisColor:
                    if candidateNode in self.graph.get(nodeWithThisColor):
                        print("Not coloring", candidateNode, "because it connects to", nodeWithThisColor)
                        shouldColor = False
                        break
                # if candidateNode in self.graph.get(node):
                #     print("Not coloring", candidateNode, "because it's adjacent to", node)
                #     continue
                if shouldColor:
                    candidateNode.color = thisColor
                    # sortedNodes.remove(candidateNode)
                    nodesWithThisColor.append(candidateNode)
                    print("Colored", candidateNode, "with color", thisColor)

            thisColor += 1
            newSortedNodes = []
            for n in sortedNodes:
                if n not in nodesWithThisColor:
                    newSortedNodes.append(n)
            sortedNodes = newSortedNodes


def main():
    G = Graph()

    onHeader = True
    # f = open('field-exclusions.CAMPAIGN_PERFORMANCE_REPORT.csv', 'r')
    f = open('example-graph.csv', 'r')
    lines = []
    for line in f:
        # skip header line
        if onHeader:
            onHeader = False
            continue
        lines.append(line)

    nodesByFieldName = {}
    for line in lines:
        fieldName, exclusionsString = line.strip().split(',')
        exclusions = exclusionsString.split(';')

        node = Node(fieldName)
        G.addNode(node)
        nodesByFieldName[fieldName] = node

    for line in lines:
        fieldName, exclusionsString = line.strip().split(',')
        exclusions = exclusionsString.split(';')

        for exclusion in exclusions:
            if exclusion == '' or exclusion is None:
                continue
            # G.addNode(nodesByFieldName[exclusion])
            G.addEdge(nodesByFieldName[fieldName], nodesByFieldName[exclusion])

    print("Graph:")
    pprint(G.graph)
    G.color()
    pprint(G.graph)


    return G

if __name__ == '__main__':
    main()

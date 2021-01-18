---
title: Topological Sort Template
date: 2021-01-18 11:17:00
categories: [Leetcode]
tags: [leetcode,dfs,topologial sort,algorithm]
---
<!--more-->

## Topological Sort template
Topological Sort template based on DFS in Java


```java
int UNKNOWN = 0, VISITING = 1, VISITED = 2;

List<Integer>> graph; // build graph
int[] status = new int[n]; // init as UNKNOWN state
List<Integer> resList = new LinkedList<>(); // reversed order

int[] topologicalSort() {        
    for (int i = 0; i < n; i++) {
        // if detected cycle, return empty array
        if (!dfs(graph, i, status)) return new int[0];
    }

    // construct final res with reverse order
    for (int i = n-1; i >= 0; i--) {
        res[n-i-1] = resList.get(i);
    }
    return res;
}

// return false if there is a cycle
boolean dfs(int cur) {
    if (status[cur] == VISITING) return false;
    if (status[cur] == VISITED) return true;
    
    status[cur] = VISITING;
    for (var adj: graph.getOrDefault(cur, new ArrayList<>())) {
        if (!dfs(graph, adj, status)) {
            return false;
        }
    }
    status[cur] = VISITED;
    resList.add(cur); // add when fully visited
    return true;
}
```
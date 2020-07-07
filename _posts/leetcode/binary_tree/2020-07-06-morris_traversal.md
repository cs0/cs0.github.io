---
title: Morris Traversal
date: 2020-07-06 20:10:00
categories: [LeetCode]
tags: [leetcode,algorithm,bianry tree,medium]
---
Morris Traversal can traverse binary tree without recursion and stack, i.e. `O(1)` space usage.
In this article, we discuss the inorder Morris Traversal. 
<!--more-->
## Introduction

**Hint: How to go back to the root after processing the left subtree?**
A: we need a temp link from the root's predecessor to the root, and we need a mechanism to remove it. 

**Hint: How to find a predecessor of a node in a binary tree?**
A: 
```java
/* Find the inorder predecessor of current */
pre = current.left; 
while (pre.right != null) 
    pre = pre.right; 
```

## Algorithm

### pseudo-code
```python
current = root
while current is not null
    if not exists current.left
        visit(current)
        current = current.right # simply go to the right subtree
    else
        predecessor = findPredecessor(current)
        if not exists predecessor.right
            predecessor.right = current
            current = current.left
        else 
            predecessor.right = null
            visit(current)
            current = current.right
``` 

### Explain
1. When you find the `predecessor` of the `current`, you add a link from the `predecessor` to the `current`, by `predecessor.right=current` (link 9).

Note here you use the right child, which makes it easy for you to later on(line 5). 

2. The link will be used to go back to the root node of the left subtree(line 5).

3. The previous link causes a cycle, and we avoid the cycle in the `findPredecessor` function(line 7) by:

```java
/* Find the inorder predecessor of current */
pre = current.left; 
while (pre.right != null && pre.right != current) 
    pre = pre.right; 
```

4. If predecessor.right is not null, it means the left sub-tree is processed. And you should process the right subtree now (`current = current.right`, line 14).

5. How to process the left subtree? Iteratively go to the left most node(line 10), until the node does not have a left node, i.e. leaf (line 3). 

6. When you go back to the root after process the left subtree, it is time to remove the link from the `predecessor` to the `current`.

7. When to `visit(current)` decide whether it is an inorder, or preorder. 

### Code
Java code:
```java
/* Function to traverse a binary tree without recursion and  
    without stack */
void MorrisTraversal(tNode root) 
{ 
    tNode current, pre; 

    if (root == null) 
        return; 

    current = root; 
    while (current != null) { 
        if (current.left == null) { // base case
            System.out.print(current.data + " "); 
            current = current.right; // this can 1. go to the right subtree, and 2. move to the root after processing the left subtree
        } 
        else { 
            /* Find the inorder predecessor of current */
            pre = current.left; 
            while (pre.right != null && pre.right != current) 
                pre = pre.right; 

            /* Make current as right child of its inorder predecessor */
            if (pre.right == null) { 
                pre.right = current; 
                current = current.left; // keep going left
            } 

            /* Revert the changes made in the 'if' part to restore the  
                original tree i.e., fix the right child of predecessor*/
            else { 
                pre.right = null; // restore
                System.out.print(current.data + " "); 
                current = current.right; 
            } /* End of if condition pre->right == NULL */

        } /* End of if condition current->left == NULL*/

    } /* End of while */
} 
```

### Time Complexity
[How is the time complexity of Morris Traversal o(n)?](https://stackoverflow.com/questions/6478063/how-is-the-time-complexity-of-morris-traversal-on)

## Leetcode
[94. Binary Tree Inorder Traversal](https://leetcode.com/problems/binary-tree-inorder-traversal/)
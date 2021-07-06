---
title: Morris Traversal
date: 2020-07-06 20:10:00
categories: [LeetCode]
tags: [leetcode,algorithm,binary tree,medium]
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
Here we present an easy to remember code structure
```java
/*
key idea is to find the predecessor, and link predecessor to the current node before going down further so we can come back
*/
    public List<Integer> inorderTraversal(TreeNode root) {
        List<Integer> res = new ArrayList<>();
        if (root == null) return res;
       
        TreeNode cur = root; // copy the root reference so we do not operate on the root
       
        while (cur != null) {
            // first find the predecessor
            // which is the right most node of its left child
            TreeNode pre = cur.left;
            while (pre != null && pre.right != null && pre.right != cur) {
                // pre != null to prevent that root.left is null
                // pre.right != null so that we ends at the right most node
                // pre.right != node so we do not enter the loop we created, and we process this later
                pre = pre.right;
            }
           
            if (pre == null) {
                // Case 1: no predecessor, which means cur.left is null
                // it is time to process the cur
                res.add(cur.val);
                // then move to right
                // note right can be null, and we will break the while loop
                cur = cur.right;
            } else if (pre.right == null) {
                // Case 2: pre.right is not yet linked to cur
                // link pre.right to cur, so we can go back to cur after processing the left child-tree
                pre.right = cur;
                // process the left child-tree
                // cur.left can not be null, otherwise pre will be null, and is processed in Case 1
                cur = cur.left;
            } else {
                // Case 3: pre.right is linked to cur
                // which means the left child-tree is processed
                // we go back to cur (done automatically, as cur's pre is pointing to cur, and we are at cur now)
                // and process cur
                res.add(cur.val);
                // reset pre link to recover the original tree structure
                pre.right = null;
                // process the right subtree
                cur = cur.right;
            }
        }
       
        return res;
    }
```

### Alternative Code, Just for reference
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
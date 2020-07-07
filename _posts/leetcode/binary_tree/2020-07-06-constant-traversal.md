---
title: Constant Space Traversal in Binary Tree with Parent Pointer
date: 2020-07-07 0:17:00
categories: [Leetcode]
tags: [leetcode,elements of programming interviews,binary tree]
---
1. The binary tree node has parent field.
2. inorder traversal with constant space.
<!--more-->

## Introduction
9.10 in Page 141

## Idea
Walk through an inorder traversal, you need to record the direction: 
1. are you going down or going up? 
2. When you go up, are you from the left subtree or the right one?

Here we use a `pre` to record the previous node, so we know the direction.

```java
public static List<Integer> inorderTraversal(BinaryTree<Integer> tree) {
    BinaryTree<Integer> pre = null, cur = tree;
    List<Integer> res = new ArrayList<>();

    while (cur != null) {
        // if going downv(we came down to cur from pre)
        if (cur.parent == pre) {
            // 
            if (cur.left != null) {
                pre = cur;
                cur = cur.left; // keep going left
            } else {
                result.add(cur.data);
                // done with left, so go right if right is not empty
                pre = cur;
                cur = (cur.right != null) 
                        ? cur.right 
                        : cur.parent; // otherwise, go up
            }
        // else if moving up from left subtree
        } else if (cur.left == pre) {
            result.add(cur.data);
            // done with left, so go right if right is not empty
            pre = cur;
            cur = (cur.right != null)
                    ? cur.right
                    : cur.parent; // otherwise, go up
        // else: done with both children, so move up
        } else {
            pre = cur;
            cur = cur.parent;
        }
    }
}
```

Note, the `pre = cur` is repeated multiple times, so we can use a `next` reference to avoid the repeat.

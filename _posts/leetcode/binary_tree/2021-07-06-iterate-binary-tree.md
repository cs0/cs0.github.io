---
title: How to iterate the binary tree without recursion correctly
date: 2021-07-06 15:36:00
categories: [LeetCode]
tags: [leetcode,algorithm,binary tree,medium]
---
How to iterate binary tree without using recursion for in-order, pre-order, and post-order. 
We should use stack, but it is not as straightforward as using recursion.
<!--more-->
## Keys To Remeber
- Simulation. The idea is to simulate the process. 

- pre-order and post-order are similar. pre-order is `root, left, right`, post-order is `left, right, root`, whose reverse is `root, right, left`. Post-order should use property, otherwise it will be hard to write it correctly. 

## In-order
```java
    public List < Integer > inorderTraversal(TreeNode root) {
        List < Integer > res = new ArrayList < > ();
        Stack < TreeNode > stack = new Stack < > ();
        TreeNode curr = root;
        while (curr != null || !stack.isEmpty()) {
            // keep pushing into stack until reach the root
            while (curr != null) {
                stack.push(curr);
                curr = curr.left;
            }
            // process the root
            curr = stack.pop();
            res.add(curr.val);
            // go the the right subtree
            curr = curr.right;
        }
        return res;
    }
```

## Pre-order
```java
    /*
    use stack, the idea is push first and visit last
    */
    public List<Integer> preorderTraversal(TreeNode root) {
        List<Integer> res = new ArrayList<>();
        if (root == null) return res;
       
        Deque<TreeNode> stack = new ArrayDeque<>();
        stack.push(root);
       
        while (!stack.isEmpty()) {
            TreeNode node = stack.pop();
            res.add(node.val);
            // need to check null
            if (node.right != null) stack.push(node.right);
            if (node.left != null) stack.push(node.left);
        }
       
        return res;
    }
```

## Post-order
```java
    /*
    using stack, hard to simulate the postorder, so we simulate the reversePostorder
    postorder:        [left, right, root]
    reversePostorder: [root, right, left]
   
    recall preorder:  [root, left, right], there is a diff between reversePostorder
    */
    public List<Integer> postorderTraversal(TreeNode root) {
        List<Integer> res = new ArrayList<>();
        if (root == null) return res;
       
        Deque<TreeNode> stack = new ArrayDeque<>();
        stack.push(root);
       
        while (!stack.isEmpty()) {
            TreeNode node = stack.pop();
            res.add(node.val); // remember it is reversed result
           
            if (node.left != null) {
                stack.push(node.left);
            }
           
            if (node.right != null) {
                stack.push(node.right);
            }
        }
       
        Collections.reverse(res);
        return res;
    }
```
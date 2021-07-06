---
title: How to write quick sort correctly
date: 2021-07-05 20:49:00
categories: [Leetcode]
tags: [leetcode,sort,quick sort,algorithm]
---
How to write quick sort correctly?
How to remember quick sort during an interview?

<!--more-->
## Three regions
Partition function is the key point. And in the textbook there are three regions for loop invariant, but it is relative hard to remember the three(two) regions: 
- `[start, i-1]` are elements less than or equal to pivot
- `[i, j]` are elements greater than pivot
- `(j, end]` are elements not scanned 

## One loop invariant
The above `i` is a bound, while the `j` is a iteration index.
It can be converted to just one loop invariant:
```
The scanned elements after `bound` are greater than pivot. 
```
So we just need to maintain one bound variable. 
Based on above we have the following partition function. 

```java

// return the partition index
private int partition(int[] nums, int start, int end) {
    int pivot = rand.nextInt(end-start+1) + start;
    swap(nums, pivot, end);
    int bound = start - 1;
    // loop invariant: scanned elements after bound are always greater than pivotNum
    for (int i = start; i < end; i++) { // i is only a scanning index
        if (nums[i] <= nums[end]) {
            swap(nums, ++bound, i);
        }
    }
    // should be ++bound here? think of an extreme case where pivot is the smallest num
    swap(nums, ++bound, end);
    return bound;
}

```

To put everything together
```java
class Solution {
    Random rand = new Random();
    public int[] sortArray(int[] nums) {
        int[] arr = Arrays.copyOf(nums, nums.length);
        // quick sort uses inclusive end
        quickSort(arr, 0, arr.length-1);
        return arr;
    }
    
    private void quickSort(int[] nums, int start, int end) {
        // base case
        if (end-start < 1) return;
        
        int partitionIndex = partition(nums, start, end);
        
        // divide and conquer
        quickSort(nums, start, partitionIndex-1);
        quickSort(nums, partitionIndex+1, end);
    }
    
    // return the partition index
    private int partition(int[] nums, int start, int end) {
        int pivot = rand.nextInt(end-start+1) + start;
        swap(nums, pivot, end);
        int bound = start - 1;
        // loop invariant: scanned elements after bound are always greater than pivotNum
        for (int i = start; i < end; i++) { // i is only a scanning index
            if (nums[i] <= nums[end]) {
                swap(nums, ++bound, i);
            }
        }
        // should be ++bound here? think of an extreme case where pivot is the smallest num
        swap(nums, ++bound, end);
        return bound;
    }
    
    private void swap(int[] nums, int i, int j) {
        int temp = nums[i];
        nums[i] = nums[j];
        nums[j] = temp;
    }
}
```
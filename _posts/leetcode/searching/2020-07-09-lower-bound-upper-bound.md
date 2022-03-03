---
title: Log(n) time to find lower bound and upper bound
date: 2020-07-07 11:17:00
categories: [Leetcode]
tags: [leetcode,elements of programming interviews,searching,binary search,algorithm]
---
For binary search of repeated elements in sorted array, the `binarySearch` in Java can not guarantee that the returned value is the first or last. So it will return a random index of these repeated elements. 

Here we discuss how to return the first or last of such repeated element. 
<!--more-->
## How to write binary search correctly?
[How to write binary search correctly](https://zhu45.org/posts/2018/Jan/12/how-to-write-binary-search-correctly/)

[Writing correct code](https://reprog.wordpress.com/2010/04/25/writing-correct-code-part-1-invariants-binary-search-part-4a/)

[Ultimate Binary Search Template](https://leetcode.com/discuss/general-discussion/786126/python-powerful-ultimate-binary-search-template-solved-many-problems)

## Ultimate Template

Suppose we have a search space. It could be an array, a range, etc. Usually it's sorted in ascending order. For most tasks, we can transform the requirement into the following generalized form:

**Minimize k , s.t. condition(k) is True**

The following code is the most generalized binary search template:

```python
def binary_search(array) -> int:
    def condition(value) -> bool:
        pass

    left, right = min(search_space), max(search_space) # could be [0, n], [1, n] etc. Depends on problem
    while left < right:
        mid = left + (right - left) // 2
        if condition(mid):
            right = mid
        else:
            left = mid + 1
    return left
```
Things to remember:
- boundary: left inclusive, right exclusive;
- if condition is satisfied, right = mid;
- return left. Remember this: after exiting the while loop, **left is the minimal k​ satisfying the condition function**;


**Before you continue, the following is just to apply the template. There is no need to read if you understand the template.**
## Array Implementation
This part is borrowed from [Stack Overflow: Implementation of C lower_bound](https://stackoverflow.com/questions/6443569/implementation-of-c-lower-bound).

`lower_bound` is almost like doing a usual binary search, except:

1. If the element isn't found, you return your current place in the search, rather than returning some null value.
2. If the element is found, you search leftward until you find a non-matching element. Then you return a pointer/iterator to the first matching element.

Note that here high index is set to n instead of n - 1. These functions can return an index which is one beyond the bounds of the array. I.e., it will return the size of the array if the search key is not found and it is greater than all the array elements.

The following can be considered as Java implementation. 

```java
public int bs_upper_bound(int[] a, int n, int x) {
    int l = 0;
    int h = n; // Not n - 1, as we use exclusive end
    while (l < h) {
        int mid =  l + (h - l) / 2;
        if (x > a[mid]) {
            l = mid + 1;
        } else if (x == a[mid]) {
            l = mid + 1; // same as above, but here we duplicate it to emphasize
        } else {
            h = mid;
        }
    }
    return l;
}
```

```java
public int bs_lower_bound(int[] a, int n, int x) {
    int l = 0;
    int h = n; // Not n - 1
    while (l < h) {
        int mid =  l + (h - l) / 2;
        if (x > a[mid]) {
            l = mid + 1;
        } else if (x == a[mid]) {
            h = mid; // same as above, but here we duplicate it to emphasize
        } else {
            h = mid;
        }
    }
    // loop ends when left == right
    return l;
}
```
Based on [How to write binary search correctly](https://zhu45.org/posts/2018/Jan/12/how-to-write-binary-search-correctly/), in the above implementation, we use the region `[l, h)` as the search region. The benefit of using such region is when the loop ends, `l == h`. So it does not matter whether you return `l` or `h`. And we randomly return `l` here. 

Then Revisit the definition of `lower_bound` and `upper_bound`:
1. `lower_bound`: the first element greater than or equal to given input. So if the input is less than the first element, return `0`. If the input is greater than the last element (with index `n-1`), return `n`.
2. `upper_bound`: the first element greater than the given input. So if the input is less than the first element, return `0`. And if greater than the last element, return `n`.

We observed for `[low, high)` invariant, `lower_bound` and `upper_bound` basic logic is the same as `binary_search` except:
1. the logic when `array[mid] == x`, instead of return directly, `lower_bound` and `upper_bound` have different search strategy.
2. return `l` at the very end instead of `-1` to indicate a not found. 

The actual C++ implementation works for all containers:
```cpp
template <class ForwardIterator, class T> ForwardIterator lower_bound (ForwardIterator first, ForwardIterator last, const T& val) {
  ForwardIterator it;
  iterator_traits<ForwardIterator>::difference_type count, step;
  count = distance(first,last);
  while (count>0)
  {
    it = first; step=count/2; advance (it,step);
    if (*it<val) {                 // or: if (comp(*it,val)), for version (2)
      first=++it;
      count-=step+1;
    }
    else count=step;
  }
  return first;
}
```

## Java List Implementation
Here we give a simplified implementation that can only work when the key exists. 
```java
/**
 * Binary search k in List A. 
 * If k occurs multiple times, return the first index.
 * If k does not exist, return -1
 * @param A the List
 * @param k the number
 */
public int lowerBound(List<Integer> A, int k) {
    // A.subList(left, right) is the candidate set
    int left = 0, right = A.size() - 1, result = -1;
    while (left < right) {
        int mid = left + ((right - left)/2);
        if (A.get(mid) > k) {
            right = mid;
        } else if (A.get(mid) == k) {
            result = mid;
            // keep searching if this is the first occurrence
            right = mid;
        } else { // A.get(mid) < k
            left = mid + 1
        }
    }
    return result;
}
```

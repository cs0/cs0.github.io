---
title: Log(n) time to find lower bound and upper bound
date: 2020-07-07 11:17:00
categories: [Leetcode]
tags: [leetcode,elements of programming interviews,searching,binary search,algorithm]
---
For binary search of repeated elements in sorted array, the `binarySearch` in Java can not guarantee that the returned value is the first or last. So it will return a random index of these repeated elements. 

Here we discuss how to return the first or last of such repeated element. 
<!--more-->

## C++ Implementation
This part is borrowed from [Stack Overflow: Implementation of C lower_bound](https://stackoverflow.com/questions/6443569/implementation-of-c-lower-bound).

`lower_bound` is almost like doing a usual binary search, except:

1. If the element isn't found, you return your current place in the search, rather than returning some null value.
2. If the element is found, you search leftward until you find a non-matching element. Then you return a pointer/iterator to the first matching element.

Note that here high index is set to n instead of n - 1. These functions can return an index which is one beyond the bounds of the array. I.e., it will return the size of the array if the search key is not found and it is greater than all the array elements.

```cpp
int bs_upper_bound(int a[], int n, int x) {
    int l = 0;
    int h = n; // Not n - 1
    while (l < h) {
        int mid =  l + (h - l) / 2;
        if (x >= a[mid]) {
            l = mid + 1;
        } else {
            h = mid;
        }
    }
    return l;
}

int bs_lower_bound(int a[], int n, int x) {
    int l = 0;
    int h = n; // Not n - 1
    while (l < h) {
        int mid =  l + (h - l) / 2;
        if (x <= a[mid]) {
            h = mid;
        } else {
            l = mid + 1;
        }
    }
    return l;
}
```

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

## Java Implementation
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
    // A.subList(left, right+1) is the candidate set
    int left = 0, right = A.size() - 1, result = -1;
    while (left <= right) {
        int mid = left + ((right - left)/2);
        if (A.get(mid) > k) {
            right = mid - 1;
        } else if (A.get(mid) == k) {
            result = mid;
            // keep searching if this is the first occurrence
            right = mid - 1;
        } else { // A.get(mid) < k
            left = mid + 1
        }
    }
    return result;
}
```

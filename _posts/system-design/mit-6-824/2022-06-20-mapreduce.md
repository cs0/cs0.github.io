---
title: MapReduce 
date: 2022-06-20 11:17:00
categories: [System Design]
tags: [system design, paper]
---
<!--more-->

## Computation Model

MapReduce is a programming model and associated implementation for processing and generating large data sets.

The computation takes a set of input key/value pairs, and produces a set of output key/value pairs.
Users specify a `map` function that processes a key/value pair to generate a set of **intermediate** key/value pairs, and a `reduce` function that **merges** all intermediate values associated with the same intermediate key.

The input keys and values are drawn from a different domain than the output keys and values. `Map` changes the key/value domain. The intermediate keys and values are from the same domain as the output keys and values.

The user specified map and reduce operations allow us to 1) parallelize large computations easily and 2) to use re-execution as the primary mechanism for fault tolerance. 

## Implementation Details

Machine failures are common. Storage is provided by inexpensive IDE disks attached **directly** to individual machines. Distributed file system is used to manage the data stored on these disks. The file system uses **replication** to provide availability and reliability on top of unreliable hardware. Users submit jobs to a scheduling system. Each job consists of a set of tasks, and is mapped by the scheduler to a set of available machines within a cluster. 

`Map` automatically partition the input data into a set of *M splits*. `Reduce` invocations are distributed by partitioning the intermediate key space into `R` pieces using a partition function. The number of partitions (`R`) and the partitioning function are specified by the user. 

## Questions

- Will reducer start in parallel with mapper, or it needs to wait for all map to finish first. 

Based on [this SO](https://stackoverflow.com/questions/13373586/does-reduce-in-mapreduce-run-right-away-or-wait-for-map-to-complete#:~:text=In%20a%20MapReduce%20job%20reducers,some%20maps%20are%20still%20running%20..),
the reducer starts to copy intermediate data as soon as they are available, but the execution will wait for all mapper to finish. 
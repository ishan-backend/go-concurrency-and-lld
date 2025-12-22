Problem Statement
Design and implement a Multiple Level Cache Management System with N levels, say:

L1 -> L2 -> L3 .... -> Ln. Each layer will store key-value pairs of data. Both KEY and VALUE are of type String. L1 is the top layer with most priority. LN is the last layer with the least priority. You are given following details about the system:

1. The number of levels of cache.
2. The capacity of each layer, i.e. number of elements that can be stored.
3. Read time of each layer.
4. Write time of each layer.

This Cache system should be able to perform following operations:

1. READ KEY Operation:
* Read will start from L1 level. 
    * If Key is found at this layer then this value will be returned.
    * If Key is not found at this layer then try to read from next layer. 
    * Keep doing this until the value of KEY is found at any level, or the last level has been reached. 
* If the value of KEY is found at any level, then this VALUE should also be written into all previous levels of cache which have higher priority that this level. 
* Total Read time is the sum of Read times of all levels where READ operation was attempted and Write time of all levels where WRITE operation was attempted. 
* You have to print the VALUE of KEY, and the total TIME taken to read it.

2. WRITE KEY Operation:
* Write will start from Ln level. To have strong consistency
    * Write the value of KEY at this level and all the levels above it.(Ln-1,Ln-2,...L1)
    * If at any level, value of KEY is same as given VALUE then don't write again and return. 
* Total Write time is the sum of Write times of all levels where WRITE operation was attempted and Read time of all levels where READ operation was attempted. 
* You have to print the total TIME taken to write this KEY-VALUE.

The eviction is happening using LRU Algorithm.

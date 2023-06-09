- sometimes we want to define additional functionalities not on one type but a collection of types. e.g. document model (list, para, etc) pe printing functionalitiy define krni hai.
- we want to have new functionality separate, since we don't want to make a lot of change in hierarchy
- the visitor is a pattern where some component of visitor is allowed to visit or traverse the entire
hierarchy of types.
- you propagate just a single method called Accept() throughout the entire
hierarchy, then you use that single method to add additional functionality, whether it's functionality for
traversal or just something else, some sort of order processing of a set of related types.


1. Intrusive Visitor Pattern
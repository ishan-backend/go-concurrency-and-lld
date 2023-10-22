* Adapter is a structural design pattern that allows objects with incompatible interfaces to collaborate.
* Give us the interface we require from the interface we have.
* Adapter = construct which adapts an existing int X, to conform to required interface Y.

https://refactoring.guru/design-patterns/adapter

**A complex scenario why you want to build an adapter in first place**:
1. API for rendering graphical objects (all images are vectored i.e. built in form of lines)
2. We are consuming this API from external system - we have a function for making a new rectangle

3. You can go ahead with adapter caching as well optionally, if you dont want to recompute same results for same request.

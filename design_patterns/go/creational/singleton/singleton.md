For some components, it makes sense to have only one initialisation in a distributed system.
1. Database repository (Redis, Postgres) - its expensive for computational resource (creating multiple times) & memory (multiple tomes created)
2. Object Factory - if a factory is stateless i.e. no fields, then there is no point having multiple instances of this factory.

e.g. construction call is expensive i.e. takes more time and consumes more resource. We want to give the same instance.
We prevent anyone from adding additional copies of an object.
Need to take care of this using lazy instantiation (use it only when somebody wants to use it).



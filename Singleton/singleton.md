- For some components it makes sense to only have one instance of it in the system.
e.g Database repository - reading all data into memory after initialising the app.
(Construction call is expensive. We want to do it only once.)

e.g. Object Factory - if factory is stateless, i.e. it doesnot have any variable. Then it does not make sense to have more than one instance of the object factory.
Reuse the instance you already have.
(We want to give everyone the same instance.)
(Want to prevent anyone from adding any additional instances.)

- Want to take care of lazy instantiation.

- Conclusion:
SINGLETON :-  A component which is instantiated only once.
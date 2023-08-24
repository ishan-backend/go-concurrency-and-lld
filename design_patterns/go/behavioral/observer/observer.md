##### Motivation

- Sometimes we need to be informed when certain things happen
  1. Object's field changes
  2. Object does something 
  3. Some external event occurs

- For point 1, (You can do polling every 100ms to figure this out) but no, we would want our object to tell us when it changed.
  We want to listen to events & notified when they occur.

- Two entities: Observer & Observable
- **Observer**: is an object which wishes to be informed about changes happening in the system.
- **Observable**: is entity generating the events. So it will have some subscribers listening to it.


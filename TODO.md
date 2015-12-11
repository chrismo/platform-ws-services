- [ ] README instructions
- [ ] Existing TODOs in code
- [ ] Listener is too generic a name. It isn't just passing events to another
      thing about alerts, it's acting on them directly.
- [ ] Why are alerts set to expire in 5 minutes?
- [ ] Why are some routes declared twice (with/without ending slash)?
      except on alerts and groups? Convenience presumably? (main.go)
- [ ] 'incidents' in README - change to alert?
- [ ] curious if the listener could be overloaded performance-wise w/
      (100ms timeout) ... or does the channel handle multiples from the alerts
     handler -- buffered? no - unbuffered, so "By default, sends and receives
     block until the other side is ready." ... I guess if Redis bogged down,
     the channel could block back to the alert handler, and perhaps subsequent
     calls to the endpoint could hang. Wonder what the connection timeout is
     to Redis. Anyway, be fun to play with.
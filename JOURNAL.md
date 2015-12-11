# Overview

Being on the receiving side of a code challenge, I like to get multiple commits
vs. squashed commits, I like READMEs telling me about next steps, challenges,
etc, because I'm as interested in the thought process as I am the final product.

Since Go is new for me and I'm sure my submission will reflect that, I thought
I'd go an extra step and journal my brain through this process. Also, I'm not
likely to spend a lot of uninterrupted time on it given the holidays
approaching.

### Dec 10 2015

Grabbed the exercise last night, was quickly reminded that Go was one of the
pieces of tech in play, and decided, sure, let's learn some Go. How bad could
it be? :)

Since I'm a long time fat IDE person (Eclipse, RubyMine), I naturally decided
to pick up Atom and try that out too. WHAT COULD GO WRONG.

Got everything installed and at least the go-plus Atom package happy with
GOPATH for now, though I read something implying GOPATH fu is a pain. Prolly
similar pain over things like CLASSPATH for Java. Every lang has its config
issues.

Browsed the .go files, seems fairly straightforward. Ah right, no exceptions, so
I get to sort of put [Joel Spolsky's
ideas](http://www.joelonsoftware.com/items/2003/10/13.html) to the test for
myself, since I guess I've never worked in an exception-less language since
before Delphi.

I saw a couple of refactor TODOs - be good to address, though I should start
with the stated goal first and see how long it takes me tackle just that.

I'm not seeing tests, I don't think? I see some things with Test in the name of
the constants, but these don't seem like unit tests? I should read up on TDD
stuff in Go.

---

Being a go n00b, trying to get the older versions of RethinkDB and (especially)
gorethink installed is leaving a little hair on my desk.

`go get` apparently has no version support for either pre-built pkg or
tag/branch from github?

Settled on `go get -u github.com/dancannon/gorethink` then cd to the $GOPATH src
and `git co v1.0.0` and `go install` from there. We'll see.

---

Got through README with a success seed, run and curl. The output is slightly
different than the README had, so I've updated that.

---

Digging into testing more, seems using build tags to create seed data is a bit
wonky to me. Maybe it's a community/golang thing. "wonky" => not a clear
declaration of intent. Seeding a dev database isn't a test.

And, confirmed `go test` says no tests found. Awesome! :|

---

In the README, I changed "deployment type" to "deployment" in the sentence
describing what a Check is because as I understand it so far, each deployment
may have many alerts/checks/incidents (presuming these are synonymous so far)? A
deployment type says to me that we record different categories of deployments,
but I see nothing other than instances of deployments so far.

If those three terms are synonymous, that could use some unification.

---

Line continuation chars inside json for alert in README need to go, at least
on bash/OS X. In for realz, I'd ask around, PR that, but...

Ah, alerts and checks have different routes.

Alerter wraps Redis. Passed into deployment handler and listener, listener
passed into alert handler. GetDeployment retrieves all active alerts from redis,
namespaced to the deployment. Listener will set or resolve alerts. Ok - that
makes sense.

---

Ok, been digging around the code flow from alert handler through alerter and
listener into Redis, including a crash course in channels. Seems
straightforward. I suppose now from the hint to have a "channel created between
listener and notifier" means polling Redis or pub/sub is uncouth -- which makes
sense as well. Why go through Redis if we can communicate directly? Dunno. (Why
store these at all in Redis, esp. with 5 min expiry?)

Also I suspect I jumped the gun on presuming deployments were instances instead
of types as the README said before.

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

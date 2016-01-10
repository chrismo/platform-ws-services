# Overview

Being on the receiving side of a code challenge, I like to get multiple commits
vs. squashed commits, I like READMEs telling me about next steps, challenges,
etc, because I'm as interested in the thought process as I am the final product.

Since Go is new for me and I'm sure my submission will reflect that, I thought
I'd go an extra step and journal my brain through this process. Also, I'm not
likely to spend a lot of uninterrupted time on it given the holidays
approaching.

### 2015 Dec 10

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

### 2015 Dec 11

Fresh start. Read over a bit of Go docs and went through the Tour examples to
help me understand some more basics, cuz reading is much easier than writing.
And I need to understand interfaces better.

So what to code next. ...

This design is bothering me. Listener, as just an interface should exist, but I
should be able to have the alerter and the notifier both implement that ...
right now it seems like passing both an alerter and a notifier into the listener
is ... no bueno. If I alter the alert handler to receive a slice of listeners,
it could loop over them all and say 'go for it' whatever you need to do. I'm
thinking this might also setup having specific notifier types for each
notification service (PagerDuty, Slack, etc.) which can then be wired up once in
main, and a generic notifier class that then is also conceptually coupled to
every notification service won't have to exist.

In real life, if there's no time for this refactoring, then we probably do the
faster thing of making Listener accept both an Alerter and a Notifier and work
to circle back on the refactoring.

I prefer the refactoring first, then the addition should just slide right into
place.

Trick now is, how do I test this to ensure refactoring don't break stuff, and
I'm still a n00b.

---

Woo! Following the time honored hack/read/run/compile/fix/read/guess cycle I
eventually moved the current Listener behind an IListener interface (nod to
Delphi) and created a passing unit test for handle_alerts.go. Whew. Such curve.
So learning.

### 2015 Dec 12

Little bit yesterday and today picked away at times on some integration tests
with Ruby. I need/want some backup on further refactoring work from the outside,
to help ensure I'm not breaking existing functionality.

I wrote it in Ruby since I'm familiar with it and seems a good job for a
scripting language to do. Also, as commented in the script, I decided to shell
out to curl instead of other native Ruby options because when I'm building out
APIs I like to include an exploratory CLI for clients that dump actual curl
commands to the console to more generically demonstrate how to use the API. The
script in its current form isn't fully exploratory, but covers some integration
basics while doing some refactoring.

### 2015 Dec 13

The current Listener is mostly code around the SensuCheck, and should probably
become a model. The alerts handler should be able to get a SensuCheck instance
go from the JSON, then just pass it to all listeners. Listener shouldn't be
parsing from JSON. Not sure why that's call a check at all, since it's really an
alert.

---

(Sheesh go can be bossy ... "Hey, comment this exported thing." "Ok, now do it
right, moron").

---

Starter Alert model going. Time to commit.

---

Listener is gutted and all but gone. SensuCheck still hanging around in
notifier. I probably broke the seeding with this next commit I'm about to do,
but the few unit tests and the integration test are passing.

---

Cool. Listener refactoring is done. Next step, "move existing notifier.getInfo
into Deployment"

---

Hmmm, tried that, but it doesn't make as much sense now digging in.

---

Got notifier as a Listener and trying to get it to hand off to transmitters
(PagerDuty, Slack). Ways to go, but it's late and eyeballs turning off.

### 2015 Dec 14

So, actually got it working before sleep last night. And a couple of more
touches before committing. Here's the bulk of the commit message:

First, the alert handler now takes an array of Listener interfaces and a
notifier instance is now passed to it in addition to the alerter.

I'd considered making model_pagerduty and model_slack simply implement the
listener interface, too, but then couldn't find a decent place for the given
`getInfo` method ... plus deciding to not use goroutines as a first shot ...
even as I type this I realize these are all very negotiable things, so I'll just
run in a direction until it seems suboptimal.

So with the notifier staying around, I added a Transmitter interface that the
specific models can implement and I have that basic functionality working in
notifier, plus a test.

---

TIL: `golint` is the bossy one. Atom's Go integration by default runs the
linter. The existing code here raises lots of complaints, but as it's a separate
tool, presumably its use is a team-by-team convention. I'm trying to follow it
in new code I'm writing I'll leave most of the existing stuff as-is.

---

Now working on pushing data to the Transmitter. Testing this now involves having
some fixture code in the RethinkDB. So, time to get some of that setup.

At the same time, I'm now in the midst of my first test expecting a test double
to be fully processed by the production code's channel before the test ends,
and that sort of async testing is usually wonksville. For now, I'll probably
fall back to directly testing the method called by the goroutine.

---

A decent compromise - bypassing the goroutine, but still using the channel. Got
DB setup working, wanna refactor some of this, e.g. extracting out common seed
and setup/fixture code whatnot.

and ... Done.

---

Hacking on the PagerDuty implementation. Not sure I have a handle on the best
way to uniquely ID a given alert. Not sure about the CapsuleID.

Ok ... have a working implementation, setup a free trial to verify via a unit
test. I'm hardcoding the Service Key which needs to be in the settings, also
want to flesh out the integration test to do this all the way through.

### 2015 Dec 15

I've decided to go with logging problems from the Transmitter, rather than
returning anything to the Notifier. In real life, we'd need to consider the best
way to error handle those things, Splunk notifications or whatever (yo dawg).
Also, with PagerDuty seems it would be a good thing to have actual daily test
notifications going to a sample alert, just to help make sure things haven't
broken down without any notice.

### 2015 Dec 16

Actually, yesterday changed my mind to have Transmitter return a TransmitResult
primarily for testing at this point. Production code ignores it, but it allows
some better assertions. Got the tests cleaned up, also added code to skip
PagerDuty if there's no setting configured. Need to do a proper integration test
next.

### 2015 Dec 17

The integration test has hit a snag in that it's presuming idempotency on
Deployment, Group and Check - but Check is giving me a problem. The code appears
to allow for replacement in the event of a conflict, but it's still not
evaluating as returning an inserting or replaced 'row' (or whatever rethinkdb
calls them).

So, hopefully this detour won't be long. I could just bypass by always
generating new values in my integration test, since this isn't directly
pertinent, but in real life I'd need to chase this bug, so for the sake of the
exercise I'll chase it a bit. At this early stage, it may not actually be a bug
but something I'm doing wrong.

---

Well, two unit tests for both Check and Deployment, and neither are working
properly. Have compared to sample code inside the rethinkdb lib source, and it
seems fine. Dunno what's wrong.

Dumping the resp to the console shows this:
```
2015/12/17 17:00:36 {0 1 0 0 0 0 0 0 0 0 0 0 0 0 []  [] []}
2015/12/17 17:00:36 {0 0 0 1 0 0 0 0 0 0 0 0 0 0 []  [] []}
```

So maybe the code to confirm something happened is referencing the wrong
counter.

RunWrite returns a WriteResponse record which has these ints:

```
type WriteResponse struct {
	Errors        int              `gorethink:"errors"`
	Inserted      int              `gorethink:"inserted"`
	Updated       int              `gorethink:"updated"`
	Unchanged     int              `gorethink:"unchanged"`
	Replaced      int              `gorethink:"replaced"`
  ...
```

So, `Unchanged` is the indicator on a true idempotent save, presumably it's
savvy enough to know nothing at all changed, which shouldn't be an error.

I guess either a check on `Errors` -- may not be enough to catch an upsert
problem, or also check `Unchanged` ... in real life I'd keep chasing this down
to make sure there was appropriate behavior.

---

Alrighty - everything finalized in the integration test and alerts going all
the way through my trial PagerDuty account.

Oh, still have to test group settings at integration level.

---

Group settings test added and verified. W00t.


### 2016 Jan 8

Long time no hack. Holidays hit with lots to do, now the end of the first week
back at it and haven't had my head in this game. Plus, OAuth is bleh. At least
I think that's where I left off on Slack. Either on the 17th or 18th of Dec,
I didn't journal anything then. IIRC, I couldn't get solid info in their docs
on how to authenticate properly, and couldn't get the code and test to work,
dug some more and confirmed I'd need to authenticate first.

So I'm back - but want to commit what I have first.

---

There. So ... I guess reading before, late at night, I couldn't decipher
between incoming web-hooks and the API for posting a message. The code here
had the API url already, but turns out an incoming web-hook works just fine,
no auth required, looks like a tokenized custom url. I'll adjust accordingly.

---

Ok, all tests passing ... ish. Added in integration tests, and it seems Slack
is only getting the deployment settings message, not the group settings one.
Which is odd because the code to make the AlertPackage is outside of the
specific Transmitter implementation. So why it'd work for one and not the
other, I'm not sure. Maybe Slack throttles their endpoint? Dunno, but I'm
committing this and done for tonight.


### 2016 Jan 9

Ah, classic hubris. "Maybe Slack throttles their endpoint," or maybe I didn't
test properly with a separate deployment instance when switching out group
settings and deployment settings. Yet again, a good night's sleep is sometimes
the best design tool.

It's possible there's also a bug somewhere here in either my integration
script (more likely) or the service when updating a deployment. In real life,
I'd chase this down, but for now going to move on - bunch of TODOs in the code
I want to look at.

---

Going to tackle the CurrentChecks "this is all kinds of gross" TODO. Adding
tests to cover a refactoring.

First thing is a deployment with no checks raises a db error, instead of just
an empty checks array. I'd think proper error handling here would be preventing
a deployment being added with no checks. CurrentChecks shouldn't necessarily
error out. Keeping that part as-is for now, need more tests first.

So ... I add a check in the next test, same error instance, I dump it out, and
now see we got other problemos: gorethink: Index `type` was not found on table
`alerts_test.checks`. Ahhh, when linking the seed scripts to the test setup,
the index was missed because it's only added in the seed method that adds the
seed data - that needs shoring up.

---

Well, that was annoying. Dunno what of my setup code is leftover from the
original exercise code and what's my screwup, but the code to create the index
wasn't being applied to the proper database -- plus there was no 'type' index
anyway. So, now that that's put to bed, and I understand things a bit better,
onward with the refactoring.

---

So, the logic in CurrentChecks that goes out of its way to merge Checks read in
from the database with any already loaded in memory seems superfluous.
Replacing it with a mere dump of what's from the database works against my
current tests, but given no pre-existing tests, I should experiment a bit here
in case I can think of a decent reason to do this. I'm guessing no.

So ... if Checks are in the Deployment Checks array, then those are separate
instances saved WITH the Deployment than Checks that are saved as standalone?

... Yes. Ok. Brain adjusted. RethinkDB Data Explorer is a thing I learned. Now
I understand why it's trying to merge these. I DON'T understand if this
separate storage is bad design or necessary.

---

Ok. README declaring the standalone Check model as a ''"dictionary" of all
checks' ... that makes sense as separate from any that are then attached and
separately stored with the Deployment. The README goes on to say "returned by
each deployment" which confuses me. I'd prefer more details on how the default
checks exist in combination with optional deployment-specific checks.

Funny how my mental model hadn't picked up on this yet. I've been presuming
the related []Check in Deployment was a reference to the Check table, not its
own storage. The code only started making sense once I realized that, then
realized a browser (enter rethink Data Explorer) would make it crystal clear.

So, I've added a method to Deployment called "DefaultChecks", plus a long test
detailing how the standalone ('default') Checks work in conjunction with the
deployment-specific checks.

Going to commit this, then I need to go BACK to the "this is gross" TODO and
see if there's some better idiomatic Go the TODO author was referring to,
now that I understand the merging functionality is intentional. That or the
"this is gross" author didn't understand the necessary complexity.

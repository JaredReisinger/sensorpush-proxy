# Background

> _A bit of background/history, relegated to the docs directory so as not to overwhelm the README._

The `sensorpush-proxy` project is an offshoot of some functionality I wrote in order to provide gratuitous infomation for a puppycam website. In particular, the site offered “Puppy Weather” information with to-the-minute temperature and humidity information taken from a small SensorPush sensor.

From that original site's README:

> ### “Puppy weather” information
> 
> Okay... now _this_ part of the site is a _**completely**_ gratuitous use of
> technology. It's a complete sub-system in and of itself. It comprises:
> 
> - [temperature/humidity data collection](#temperature-humidity-data-collection)
> - [bespoke data proxy](#bespoke-data-proxy)
> - [auto-updating in-browser logic](#auto-updating-in-browser-logic)
> 
> #### Temperature/humidity data collection
> 
> For various and sundry reasons, I happen to have several
> [SensorPush](https://www.sensorpush.com/) temperature and humidity sensors. They
> allow me to monitor such critical things as _“does my cigar humidor need more
> water?”_ and _“what temperature is the kitchen right now?”_. These devices
> record the temperature and humidity once every minute, and send that information
> (via a low-energy-Bluetooth–to–WiFi hub) up to SenorPush’s data servers in the
> cloud.
> 
> SensorPush also exposes a programmatic API so that one can get the raw data they
> have saved from the individual sensors. Their API is rather complex, however,
> and also requires performing proper authentication using one’s (my!) personal
> SensorPush password. Since (a) I don’t want to expose my password even within
> the source of the web page, (b) the PuppyCam site only needs the most-recent
> temperature/humidity reading, and (c) when the PuppyCam goes viral it could
> start requesting the information _millions_ of times 😉, the “most sane” thing
> to do is to run a tiny one-off web server whose job is to peridocially poke the
> SensorPush API to get the latest reading, and to also serve that information
> publically.
> 
> All of which means it’s a…
> 
> #### Bespoke data proxy
> 
> I’ve written a small program/web-server in [Go](https://golang.org/) that
> performs all of the necessary work for getting the temperature and humidity from
> the SensorPush API, and then turns around and exposes it as a simple data feed
> at _[REDACTED]_. If you visit that link,
> you’ll see some additional potential data metrics without values: altitude,
> barometric pressure, and so on. SensorPush has some newer sensors that report
> those values, but ours only includes the temperature and humidity.
> 
> This small program runs on spare hardware that’s been cobbled together over
> the years and turned into a little Linux box running my own personal private
> “cloud”. It could just as easily run “in the (real) cloud”—and in
> would/will, if the PuppyCam actually goes viral!—but it doesn’t cost
> anything (more) to run it on a machine that I already have and that is
> already exposed to the internet as half-a-dozen other virtual servers.
> 
> The source for this is in a private repository... but I may open it up in the
> future.

As that excerpt suggests at the end, “I may open it up in the future…” and the future is now!

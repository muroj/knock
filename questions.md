1. Why is it important to measure percentile latencies in production systems (e.g. p99)?

Percentiles are a better indicator of actual user experience verses more common metrics like average, minimum, or max, which can be affected by outliers. For example, a temporary network hiccup in third-party vendor software or external integration can skew the max request latency for some users. Trying to optimize your application for this (exceptional) case is normally not worth the effort. A percentile-based metric excludes outliers like this and provides a better target to optimize for.

2. Which metrics are important to track for queues? Why?

The size of the queue over a given time period. A sudden spike in queue entries could be a legitamite increase in traffic, in which case, it might be necessary to scale the system accordingly. However, it can also indicate a problem, such as a DDOS attack or an outage in a downstream processing system (i.e. consumer). Conversely, a drop in queue entries could indicate a networking problem as new requests are not being added due (i.e. a producer failure).

The processing time for each queue entry. An increase in processing time can indicate a resource constraint (or outage) on the consumer side. 

Depending on the type of queue, (e.g. a messaging queue such as RabbitMQ) the size and type of each queue message might be relevant. This information could help optimize queue configuration and system resources.
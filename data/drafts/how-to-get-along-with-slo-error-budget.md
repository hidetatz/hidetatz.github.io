How to get along with SLO and Error budget---2020-02-03 12:00:00

## Introduction

I get to know SLO and error budget when reading a book "Site Reliability Engineering"[^1]. While I am building some microservices in my company, I found how SLO and error budget are such powerful tool to develop and maintain them.
In this article, I will try to describe how SLO and error budget should work and why they such matters. I hope this article can help readers who are trying to understand and build their own SLO strategy.

## SLA

SLA is an **agreement** of service level of a system. Usually, the agreement is had between service provider and service customer.
Under SLA, service provider and customers agrees on what service level **must be** provided. Basically, they also agrees what happens when agreed service level is not met. 
This helps how customers should choose services. If the customer has a requirement about service level, then they can compare some services on their SLA.

If the SLA is not met, sometimes service provider gives a refund to the customer. They also can take other forms.
For example, SLA of AWS compute service[^2] is defined as Uptime. When Uptime is under 99.99% (and above 99.0%), 10% refund can be given.
GCP[^3], Azure[^4], Datadog[^5] and more services have SLA.

Developers or SRE don't have to get involved in defining SLA. Because SLA is more close to business and product side. Their target is usually SLO. The relation between SLO and SLA will be described later.

## SLI

SLI is a service level **indicator**. SLI is a quantitative measure which is a part of the service level. There are some common indicators.

* Request latency - How long it takes to return a response
* Error Rate - How many 5xx response is responded as a fraction
* System throughput - Typically measured as RPS (requests per second)
* Availability - A fraction of the time that a service can be used.
* Durability - How it is likely that data will be retained over a long period (e.g. 1 year). This is especially used for storage system.
* Consistency - In Microsoft Azure Cosmos DB, users can choose which consistency model they need. They provide SLA for its violation rate.[^6]

Good SLI can measure what users are interested in the service directly.
But sometimes it's difficult. In that case, using another proxy is also OK.

## SLO

SLO represents service level **objective**.
SLO usually includes SLI, which should be tracked.
For example, when we are building Storage system, and we think we want to provide high durability to customers, we can have an SLO that "99.999999% of objects won't be lost or compromised in the event of a failure over 1 year".

Do you wonder what's the difference between SLO and SLA? Actually, they are a kind of similar.
First, as described above, there will be refund or some other penalties if SLA is not met.
However, when SLO is not met, there should be no user complaints. SLA is more for users, but SLO is for developers and SREs.
Services can have both SLO and SLA. When they have both, usually SLO is not public to users while SLA is.
For example, if they have an SLA that 99.9% availability must be kept, then they can have the same SLI for their SLO. But it will be more strict value (e.g. 99.95%). In this case, SLA is a promise with users, but SLO is a target for developers. When SLO is not met, there should be no user-impact, but SLO must be always met. That's why SLO must be more strict than SLA.

The purpose is also different. SLA is for users; SLA should help users if they can choose to use the service. But SLO is for developers and SREs. SLO helps developers to prioritize their work around the service.
Let's say we are operating a service. Recently our monitoring system shows there is some delay in the service, which originally was not found. If we don't have SLO about latency, we always have to decide if the delay must be investigated or fixed. What if the latency increased by 5ms when we released a new feature? Is it a problem to be fixed? What if it is 50ms?
SLO helps this situation. We can decide what to do when we face a problem. Simply, if it violates the SLO, then stop feature development and work on fixint the problem. If SLO is still met, then keep working on feature development.
Usually, feature development and site reliability is trade-off. Typical infrastructure engineers work on only site reliability, but from Site Reliability Engineering's perspective, they should also work on feature development.
Pulling up the availability from 99.9% to 99.99% is super hard, while 99.9% can be sufficient in most cases. Having good SLO helps us to decide if we have to work on improving site reliability, or if we can keep working on feature development.

## How to implement SLO

There are some ways how to implement SLO. I just want to describe one example.

### Determine what kind of SLI is the best for the service

First, decide what should be the indicators of the level for the service.
Because SLO should include SLI, it's also determined based on users' interests. Good SLI can be a measurement what users are interested in.
Choosing good SLI requires us to understand how users most interact with the service. Choosing an SLI which users don't care about doesn't make sense.

Too many SLI makes it hard to track and keep paying attention. Too few is also not good because it usually cannot show system's health properly.

The SLI which indicates about service's correctness is recommended. For distributed database, customers usually want to know if the latest data is always returned. Sometimes it's difficult to track correctness, but it's better to consider to provide correctness as SLO.

In addition, there are some common patterns what to be chosen by the category of the service.

When the service is serving something to users, usually **availability**, **latency**, and **throughput** are chosen.
When the service is storage, **latency**, **availability**, **durability** are chosen.
When the service is batch or big data analysis platform, **throughput** (how much data can be processed) is chosen.

TODO: Write types of components

TODO: Write SLI Specification and implementation

### Decide how to measure SLIs and measure them

### Decide SLO

Deciding SLO consists of 2 parts; deciding SLI and deciding what number to set, e.g. 99.9%. We already decided SLI in above section, so then let's talk about its number.

Originally, traditional infrastructure engineers tried to keep 100% availability. Their responsibility was monitoring systems and whenever a problem happens, fixing it, despite of how important it is.
However, from site reliability's perspective, it is not the right direction we should go.

* 100% is TOO difficult as objective. Technically, multiple components simulataneous failure cannot be completely avoided. Design for failure, such as failover, server redunduncy, is not a sliver bullet.
* Because whole system consists of multiple unreliable factors (e.g. load balancer, network, users' device), users won't experience 100% availability. Even if one subsystem's availability is 100%, another 99% subsystems makes whole system's availability under 100.
* If SLO is set 100%, then implementing new features, upgrading dependency versions, applying security patches are very difficult. SLO is something which **must** be met. You cannot use time for product implovements; all you can do is always monitoring systems and being reactive. Site reliability engineering targets the balance of continuous system improvement and system reliability. SLO should not be set as 100%.
### Decide time window

### Calculate Error budgets

### After SLO is defined
* Stakeholders agreements
* Documentation
* Dashboards and monitors

---

## References

[^1]: [Site Reliability Engineering](https://landing.google.com/sre/books/)
[^2]: [Amazon Compute Service Level Agreement](https://aws.amazon.com/compute/sla/)
[^3]: [Google Cloud Platform Service Level Agreements](https://cloud.google.com/terms/sla/)
[^4]: [Microsoft Azure Service Level Agreements](https://azure.microsoft.com/en-us/support/legal/sla/)
[^5]: [Datadog Service Terms and Agreement](https://www.datadoghq.com/legal/terms/2014-12-31/)
[^6]: [SLA for Azure Cosmos DB](https://azure.microsoft.com/en-us/support/legal/sla/cosmos-db/v1_3/)

* https://landing.google.com/sre/sre-book/chapters/service-level-objectives/
* https://landing.google.com/sre/workbook/chapters/implementing-slos/
* https://landing.google.com/sre/sre-book/chapters/embracing-risk/
* https://www.usenix.org/sites/default/files/conference/protected-files/sre19amer_slides__lawson.pdf

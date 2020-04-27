How to get along with SLO and Error budget---2020-04-28 12:00:00

## Introduction

I get to know SLO and error budget when reading a book "Site Reliability Engineering"[1]. While I am building some microservices in my company, I found how SLO and error budget are such powerful tool to develop and maintain them.
In this article, I will try to describe how SLO and error budget should work and why they such matters. I hope this article can help readers who are trying to understand and build their own SLO strategy.

## SLA

SLA is an **agreement** of service level of a system. Usually, the agreement is had between service provider and service customer.
Under SLA, service provider and customers agrees on what service level **must be** provided. Basically, they also agrees what happens when agreed service level is not met. 
This helps how customers should choose services. If the customer has a requirement about service level, then they can compare some services on their SLA.

If the SLA is not met, sometimes service provider gives a refund to the customer. They also can take other forms.
For example, SLA of AWS compute service[2] is defined as Uptime. When Uptime is under 99.99% (and above 99.0%), 10% refund can be given.
GCP[3], Azure[4], Datadog[5] and more services have SLA.

Developers or SRE don't have to get involved in defining SLA. Because SLA is more close to business and product side. Their target is usually SLO. The relation between SLO and SLA will be described later.

## SLI

SLI is a service level **indicator**. SLI is a quantitative measure which is a part of the service level. There are some common indicators.

* Request latency - How long it takes to return a response
* Error Rate - How many 5xx response is responded as a fraction
* System throughput - Typically measured as RPS (requests per second)
* Availability - A fraction of the time that a service can be used.
* Durability - How it is likely that data will be retained over a long period (e.g. 1 year). This is especially used for storage system.
* Consistency - In Microsoft Azure Cosmos DB, users can choose which consistency model they need. They provide SLA for its violation rate.[6]

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

There are some ways how to implement SLO. Below is just one example.

### Determine what kind of SLI is the best for the service

First, decide what should be the indicators of the level for the service.
Because SLO should include SLI, it's also determined based on users' interests. Good SLI can be a measurement what users are interested in.
Choosing good SLI requires us to understand how users most interact with the service. Choosing an SLI which users don't care about doesn't make sense.

Too many SLI makes it hard to track and keep paying attention. Too few is also not good because it usually cannot show system's health properly.

First, the SLI which indicates about service's correctness is recommended. For distributed database, customers usually want to know if the latest data is always returned. Sometimes it's difficult to track correctness, but it's better to consider to provide correctness as SLO.

In addition, there are some common patterns what to be chosen by the category of the service.

When the service is serving something to users, usually **availability**, **latency**, and **throughput** are chosen.
When the service is storage, **latency**, **availability**, **durability** are chosen.
When the service is batch or big data analysis platform, **throughput** (how much data can be processed) is chosen.

This part is important and worth to spend time to be discussed in deciding service's SLO. Without deep consideration, we tend to create SLO just following what other teams are doing; especially when the team is not familiar with working with SLO. However, it's not a good way.
**Potentially, SLO can decide your future tasks.**
You must be an expert of your service. So, deciding SLO requires your domain knowledge.

For example, your service has WRITE operations, you might have take care of data write latency, message queue delay, message durability, etc.
If your service has only READ operations (such as master-data service), you may just take care of read latency.

When your application is called "microservice", one common mistake is just preparing SLI about application's health.
Popular ones would be "availability" and "latency". However, these SLIs are based on infrastructure's health.
If your infrastructure gets down, can your monitoring system send alert? If no, you'd better to consider adding them as your SLI.

### SLI Specification and implementation

SLI can be devided into 2 parts; **specification** and **implementation**.

#### SLI specification

__SLI specification__ is indicator to represent what matters to users. It doesn't include how to be measured.
e.g. **The latency of an API is less than 150ms**

#### SLI implementation

__SLI implementation__ describes how to measure SLI specification.

e.g.

* Latency of an API measured by load balancer's log.
* Latency of an API measured by server's log.
* Latency of an API measured by Datadog agent on server node.
* Latency of an API measured by Datadog synthetic monitoring.

To define SLI specification/implementation, first we need to understand there should be multiple measure ways of the specification. 
We can have a question like this: what's the **latency**? Should it be on server-side? On Internet-service-provider? On real user? Defining SLI implementation will help to answer the question.
There are no answers which measurement is "correct". It should be depending on what matters to your users.

### Decide how to measure SLIs and measure them

After deciding SLI specification/implementation, now we want to talk about __how to implement SLI implementation__. If your team or company is already introduced any monitoring tool such as New Relic, Nagios, Datadog, Zabbix, etc, it's a good idea to use them for SLI.
For example, Datadog[7] has features for SLI/SLO measurement.
It's of course OK to create a system to measure your SLI by yourself, but basically it's not recommended. 
If you prepare something, you need another monitoring - **the original monitoring system must be monitored** .
If you can use managed monitoring system, it will be better.

For example, in Datadog, "95 percentile of sum of HTTP server's latency measured by Datadog agent on server node" can be written (as JSON) like this:

```json
{
  "q": "sum:trace.http.server.duration.by.resource_service.95p{service:your_service_name,env:production} by {resource_name}"
}
```

More your SLI implementation is detailed, easier to implement it.

### Decide SLO

Deciding SLO consists of 2 parts; deciding SLI and deciding what number to set, e.g. 99.9%. We already decided SLI in above section, so then let's talk about its number.

Originally, traditional infrastructure engineers tried to keep 100% availability. Their responsibility was monitoring systems and whenever a problem happens, fixing it, despite of how important it is.
However, from site reliability's perspective, it is not the right direction we should go.

* 100% is TOO difficult as objective. Technically, multiple components simulataneous failure cannot be completely avoided. Design for failure, such as failover, server redunduncy, is not a sliver bullet.
* Because whole system consists of multiple unreliable factors (e.g. load balancer, network, users' device), users won't experience 100% availability. Even if one subsystem's availability is 100%, another 99% subsystems makes whole system's availability under 100.
* If SLO is set 100%, then implementing new features, upgrading dependency versions, applying security patches... become very difficult. SLO is something which **must** be met. You cannot use time for product implovements; all you can do is always monitoring systems and being **reactive** - literally. Site reliability engineering targets the balance of continuous system improvement and system reliability. SLO should not be set as 100%.

So, if you don't choose 100%, what number is good?
One common way to decide it is using past data.
Let's say we have below data collected in past 1 month:

|request count|200 response|
|:---|:---|
|1,000,000|978,857|

This data shows us that past 1 month availability is 97.8857%.

First, we set lower number than the data as SLO; e.g. 97.5%. Don't choose 98%! Too strict SLO is also bad; you can adjust them later on and if you announce your service's SLO to other teams, you basically cannot lower it.

You might think like "97% is too low! Do we need to take any actions to pull up the availability?".
Typically, it's very difficult to know what is desired number of SLO. Actually, **nobody knows it**. Product manager, data analyst, executives... won't know it. You will eventually know it after some time passed. If your SLO is too low, users may complain. If it's too high, your team members may get tired of daily operation.
Finding good SLO is not easy.

By the same method, you will be able to decide what latency you should target.

### Decide time window

One thing we need to consider in addition is **time window** - SLO's time interval.
What is time window?
Let's say we have SLO "availability, 97%" with a **2 weeks** time window.
This means, "if first 1 week's availability is under 97%, no immediate solution will be proposed. In latter 1 week, it the SLO is lifted up above 97%, it's totally fine."

Longer time window gives the team more time - while users are seeing more unavailability.

Shorter time window is more kind to users. However, you have to check SLO more often. You can make decisions more quickly.

This might affect to your project work. Basically, your time window should be longer than your team's sprint iteration. Because, you have to take some actions and create tickets for it and add into backlog.
Also, SLO can help to plan your team's headcount - you might need more experienced engineers to keep your SLO.
Google[8] recommends "4 weeks" as general-purpose time window, but it's up to your team.

### Calculate Error budgets

Finally, we are here to talk about error budget!
SLO is target percentage, error budget is 100% - SLO.

Let's use the previous data for a SLO "availability, 97%".
If we have 10,000,000 requests in 4 weeks time window, `10,000,000 * 0.03 = 300,000` so 300,000 requests can fail. This is error budget for the SLO.

After error budget calculation, monitor it. When monitoring system reports error budget getting empty, we need to take some actions; typically, like below.

* Developers must prioritize the highest to fix issues related to the SLO. If it's caused by logic bug, it must be resolved. If it's because of unstable infrastructure, they have to investigate the problem.
* Developers have to try to fix the problem until SLO is met and error budget gets surplus. During the time, no new feature requests should be accpted.
* (Optional) waiting for enough error budget are saved is good idea. If bug fix contains another problem, the SLO might get much lower.

I believe this is the most essential part in SLO work. Just deciding SLO is not enough. Without error budget, we cannot balance site reliability and feature development.
As long as SLO is met, it's not problem at all. What we want to focus on is not only keeping SLO met, but also doing what to do when SLO is not met. It's a part of **design for failure**.

### After SLO is decided

Once we decided SLO by engineer side, we need more things to be done.

#### Stakeholders agreements

SLO is targeting better development cycle, but it should not be closed in only engineer team.
Most product managers don't know what SLO is and tradeoff between site reliability and feature development.
From error budget's perspective, if it is fully consumed, developers stop feature development. It should be agreed with product managers.
Talk with them and get them understand why SLO/error budget is needed.

After getting product managers' understanding, writing documents and visualizing SLO/error budget is important. They will be described later.

#### Documentation

After SLO is decided, they should be documented and able to be found easily from product managers and other teams.
What to be documented are:

* Author, reviewers, approvers. Especially, who made the decision if it is fine from business perspective.
* Edit history. It's better to leave history how SLO has been changed.
* SLI implementation
* SLO
* Error budget calculation method and operation policy
  - When to take actions
  - What action the team will take

An example of SLO document is available by Google[9].

#### Dashboards and alerts

After SLO is decided, it should be visualized in dashboard and always visible.
What to be shown are:

* SLI timeseries graph. e.g. latency, throughput.
* SLO binary graph - OK or Not met.
* Remaining error budget time series graph.

#### Improve SLI/SLO

SLI/SLO can be improved after some actual operation.
When and why we want to change them? There are some signs:

* If error budget is consumed very often, the SLO might be too tight. Is the SLO really necessary for users?
* If users complains about your service but no alerts are shown, the SLO might be too relaxed and need to be tighter. SLO which cannot represent users' interaction would be pointless.

When relaxing a SLO, it has to get agreed with stakeholders - product managers might say it's too low and has some drawbacks to users.

Chances to review past SLO should be held periodically. When time window is 4weeks, typically once a quarter is good to see past 12 weeks SLO.
At the meeting, see past SLI graph and error budget comsumption, then talk about making it better.

## Summary

In this article, I tried to describe how to decide SLO and work with it.
Deciding SLO is not goal; goal is **doing both, site reliability and feature development**. SLO and error budget can be a game changer if you are tired of going endless journey to 100% reliability.
I hope this article can help you to decide your own SLO.

---

## References

[1]: Site Reliability Engineering [https://landing.google.com/sre/books/](https://landing.google.com/sre/books/)
[2]: Amazon Compute Service Level Agreement [https://aws.amazon.com/compute/sla/](https://aws.amazon.com/compute/sla)
[3]: Google Cloud Platform Service Level Agreements [https://cloud.google.com/terms/sla/](https://cloud.google.com/terms/sla/)
[4]: Microsoft Azure Service Level Agreements [https://azure.microsoft.com/en-us/support/legal/sla/](https://azure.microsoft.com/en-us/support/legal/sla/)
[5]: Datadog Service Terms and Agreement [https://www.datadoghq.com/legal/terms/2014-12-31/](https://www.datadoghq.com/legal/terms/2014-12-31/)
[6]: SLA for Azure Cosmos DB [https://azure.microsoft.com/en-us/support/legal/sla/cosmos-db/v1_3/](https://azure.microsoft.com/en-us/support/legal/sla/cosmos-db/v1_3/)
[7]: Track SLIs and SLOs [https://www.datadoghq.com/videos/track-sli-slo/](https://www.datadoghq.com/videos/track-sli-slo/)
[8]: Implementing SLOs [https://landing.google.com/sre/workbook/chapters/implementing-slos/](https://landing.google.com/sre/workbook/chapters/implementing-slos/)
[9]: Example SLO document [https://landing.google.com/sre/workbook/chapters/slo-document/](https://landing.google.com/sre/workbook/chapters/slo-document/)

* https://landing.google.com/sre/sre-book/chapters/service-level-objectives/
* https://landing.google.com/sre/sre-book/chapters/embracing-risk/
* https://www.usenix.org/sites/default/files/conference/protected-files/sre19amer_slides__lawson.pdf

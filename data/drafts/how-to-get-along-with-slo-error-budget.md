How to get along with SLO and Error budget---2020-02-03 12:00:00

## Introduction

I get to know SLO and error budget when reading a book "Site Reliability Engineering"[^1]. While I am building some microservices in my company, I found how SLO and error budget are such powerful tool to develop and maintain them.
Now, I want to write an article which describes XXX.
In this article, I will try to describe what SLO/SLA/SLI are and how 

## SLA

SLA is an **agreement** of service level of a system. Usually, the agreement is had between service provider and service customer.
Under SLA, service provider and customers agrees on what service level **must be** provided. Basically, they also agrees what happens when agreed service level is not met. 
This helps how customers should choose services. If the customer has a requirement about service level, then they can compare some services on their SLA.

If the SLA is not met, sometimes service provider gives a refund to the customer. They also can take other forms.
For example, SLA of AWS compute service[^2] is defined as Uptime. When Uptime is under 99.99% (and above 99.0%), 10% refund can be given.
GCP[^3], Azure[^4], Datadog[^5] and more services have SLA.

Developers or SRE doesn't have to get involved in defining SLA. Because SLA is more close to business and product side. Their target is usually SLO. The relation between SLO and SLA will be described later.

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
However, 

---

[^1]: [Site Reliability Engineering](https://landing.google.com/sre/books/)
[^2]: [Amazon Compute Service Level Agreement](https://aws.amazon.com/compute/sla/)
[^3]: [Google Cloud Platform Service Level Agreements](https://cloud.google.com/terms/sla/)
[^4]: [Microsoft Azure Service Level Agreements](https://azure.microsoft.com/en-us/support/legal/sla/)
[^5]: [Datadog Service Terms and Agreement](https://www.datadoghq.com/legal/terms/2014-12-31/)
[^6]: [SLA for Azure Cosmos DB](https://azure.microsoft.com/en-us/support/legal/sla/cosmos-db/v1_3/)

---
https://landing.google.com/sre/sre-book/chapters/service-level-objectives/
https://landing.google.com/sre/workbook/chapters/implementing-slos/
https://landing.google.com/sre/sre-book/chapters/embracing-risk/

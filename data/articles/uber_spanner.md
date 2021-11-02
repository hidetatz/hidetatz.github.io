type: input
timestamp: 2021-11-02 23:24:30
url: https://eng.uber.com/building-ubers-fulfillment-platform/
lang: en
---


* Uber's fulfillment platform handles billions of DB transactions each day
* "The platform handles millions of concurrent users and billions of trips per month across over ten thousand cities and billions of database transactions a day."
To introduce Cloud Spanner into Uber’s environment, we had to solve three main challenges:
  * How do we design our application workload that assumes NoSQL paradigms to work with a NewSQL-based architecture?
  * How do we build resilient and scalable networking architecture so that we can leverage Cloud Spanner, irrespective of where Uber’s operational regions are located (Uber On-Prem, AWS, or GCP)?
  * How do we optimize and operationalize a completely new cloud database that can handle Uber’s scale?
* Challanges:
  * Based on Ringpop, the data should be uniformly distributed, but actually it results in hotspots in some highly active jobs/supplies.
  * Ringpop based architecture trades consistency to achieve AP guarantee. Because of it, debugging production issues were really difficult. Also, some saga failures results data inconsistency and developers must manually fix it.
* By the architecture change, consistency became one of the primary evaluation criteria as well as resilliency and availability.
* Spanner: external consistency. The strictest concurrency control model where the system behaves asif all transactions were excuted sequentially, even though the system actually runs them across multiple servers.
* Create own Spanner client
  * Spanner: sessions best practices: https://cloud.google.com/spanner/docs/sessions#best_practices_when_creating_a_client_library_or_using_restrpc
  * Use gRPC client interceptor for RPC level observability and retries/timeouts for each RPC for resiliency
  * Transaction prioritization: user-facing transactions hould be prioritized over background transactions. using timeouts/retries
* LATE: original CDC solution for Cloud Spanner
  * Table: Task Table is used to store a task.
  * Tailer: fetches a task from the Task table to consume new tasks.
  * Sharding: Because a tailer works on a single shard, to distribute trailers across N tailer workers. Use a combination of Rendezvous-Hashing
* On-prem cache


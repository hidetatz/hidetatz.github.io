type: input
timestamp: 2021-10-30 23:23:34
url: https://slack.engineering/infrastructure-observability-for-changing-the-spend-curve
lang: en
---

* Data driven infrastructure decisions
* Focus on Infrastrucutre changes observability
* Observability is culture and practice for ogranizations
* CI in Slack contains a variety of tests like unit tests, integration tests, and end-to-end functional tests for a variety of codebases
* In webapp CI, a service Checkpoint, developed in Slack internally, orchestrates the complex test workflow. It also handles test failures notification, review requests, etc.
* Checkpoint can show a dashboard to analyze developer experience with views for reliability and performance of specific dimensions, like flakiness per suite, time to mergeable, and cycle time
  * https://medium.com/azimolabs/what-is-flakiness-and-how-we-deal-with-it-39b270ed5445
  * https://sourcelevel.io/blog/5-metrics-engineering-managers-can-extract-from-pull-requests
  * https://docs.velocity.codeclimate.com/en/articles/2913508-cycle-time
* In Slack they uses the dashboard to understand curve changes
* They create project hypotheses and scope potential project impact through observability through metrics, monitoring, and traces
* The story about adding circuit breakers in the CI is interesting. Because usually CI is implemented as workflow, they can experience cascading failures without it.

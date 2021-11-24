Sharing what I've read, listened, watched, etc.

## [How Dropbox Replay keeps everyone in sync - Dropbox](https://dropbox.tech/application/how-dropbox-replay-keeps-everyone-in-sync)
### 2021/11/24

* Dropbox Replay defines happend-before relationship to each message and processes them one by one
  - Using mutex, each message is received and broadcasted one by one and next message is never processed if the prev message is still in progress
* Internally the sync service must handle distributed lock? or it might be using RDB transaction for this

## [Squid game: how we load-tested Ably’s Control API | Ably Blog: Data in Motion](https://ably.com/blog/how-we-load-tested-control-api)
### 2021/11/15

* They devide load test scenarios into three;
  - Typical users
  - Power users
    - Too many traffic, but don't have malcious intent
    - To simulate this properly without facing rate limiting, they use a reverse proxy (squid) to distribute the client IP address
  - Bad users/bots
    - Send too many traffic with bad intentions

## [GitHub - shellgei/shellgei160: 書籍: シェル・ワンライナー160本ノックの情報ページ](https://github.com/shellgei/shellgei160)
### 2021/11/11

* Reading a book abount shell oneliner: https://gihyo.jp/book/2021/978-4-297-12267-6
* in: abcdefg, out: abcdbcdefg
  - `echo abcdefg | sed 's/bcd/&&/'`
* in: abcdefg, out: aefbcdg
  - `echo abcdefg | sed -E 's/(bcd)(ef)/\2\1/'`
* in: $(seq 100), out: 1 3 5 7 9 (omitted) 95 97 99
  - `seq 100 | grep "^.*[13579]$" | xargs`
    - ends with 02468.
  - `seq 100 | grep "[^02468]$" | xargs`
    - does not end with 02468.
* in: $(seq 100), out: 11 22 33 44 55 66 77 88 99
  - `seq 100 | grep -E "^(.)\1"$ | xargs`
* `grep -o`: shows matched parts only
* in: $(seq 5), out: 2 4
  - `seq 5 | awk '/[24]/' | xargs`
  - `seq 5 | awk '$1%2==0' | xargs`
  - `seq 5 | awk '$1%2==0{print $1}' | xargs`
* in: $(seq 5), out: 2 even 4 even
  - `seq 5 | awk '$1%2==0{print $1, "even"}' | xargs`
  - `seq 5 | awk '$1%2==0{print($1, "even")}' | xargs`
* in: $(seq 5), out: 1 odd 2 even 3 odd 4 even 5 odd
  - `seq 5 | awk '$1%2==0{print($1, "even")}$1%2==1{print($1, "odd")}' | xargs`
  - `seq 5 | awk '$1%2==0{print($1, "even")}$1%2{print($1, "odd")}' | xargs`
    - in awk, a condition is called "pattern", and the process is called "action"
    - in above command, "$1%2==0" is pattern and the correlated action is "{print($1, "even")}"
    - multiple condition/action can be written
* in: $(seq 5), out: 1 odd 2 even 3 odd 4 even 5 odd sum 15
  - `seq 5 | awk 'BEGIN{sum=0}$1%2==0{print $1, "even"}$1%2{print $1, "odd"}{sum+=$1}END{print "sum", sum}' | xargs`
    - BEGIN pattern matches when the awk starts to process the first line
    - END pattern matches after the awk finishes to process the last line
    - fourth action (`{sum+=$1}`) is called in every line process because there is no correlated pattern
* in: $(seq 5), out: odd 3 even 2
  - `seq 5 | awk '{print $1%2==0 ? "even" : "odd"}' | sort | uniq -c | awk '{print $2, $1}' | gsort -k2,2nr`
  - `seq 5 | awk '{print $1%2==0 ? "even" : "odd"}' | awk '{m[$1]+=1}END{for(key in m)print key, m[key]}'`
    - you can use hashmap in awk!

## [Verica - MTTR is a Misleading Metric—Now What?](https://www.verica.io/blog/mttr-is-a-misleading-metric-now-what/)
### 2021/11/08

* MTTx is trying to simplify essentially complex things too much. I totally agree with it.
* The benefit of MTTx for me looks like that we can simply tell external stakeholders "our company's MTTR is..."
* But internally, we should track SLO.

## [You Are Not Google. Software engineers go crazy for the… | by Oz Nova | Bradfield](https://blog.bradfieldcs.com/you-are-not-google-84912cf44afb)
### 2021/11/04

* You are not Google, Amazon, LinkedIn...
* Don't follow tech giants choice.
* Your problem is not theirs. Only you can understand the problem and find a solution.
* For me, Kubernetes is the same for most companies.

## [The process: How Twilio scaled its engineering structure – Increment: Teams](https://increment.com/teams/how-twilio-scaled-its-engineering-structure/)
### 2021/11/04

* To 40 members: Write things down. Describe the context, tasks, etc. Make small teams.
* To 100 members: Because of teams by expertise and skills, leaders needed to do many context-switching in a day.
  - Also the accountability was shared. However, shared accountability leads people to the lack of accountability. They changed the structure and each team (up to 10 people) has ownerships on solutions/product.
* To 1000 members: Today, a new team in Twilio get a seed investment from the company including budget, resources, and head count.
  - 150 or more teams are running their development/ops on self service platform. They don't need to spend their time on very basic stuffs, instead, they can perform things rapidly.

## [Taming Go’s Memory Usage, or How We Avoided Rewriting Our Client in Rust — Akita Software](https://www.akitasoftware.com/blog-posts/taming-gos-memory-usage-or-how-we-avoided-rewriting-our-client-in-rust)
### 2021/11/04

* The lessons learned are interesting:
  - **Reduce fixed overhead**. Go’s garbage collection ensures that you pay for each live byte with another byte of system memory. Keeping fixed overhead low will reduce resident set size.
  - **Profile allocation, not just live data**. This reveals what is making the Go garbage collector perform work, and spikes in memory usage are usually due to increased activity at those sites.
  - **Stream, don’t buffer**. It’s a common mistake to collect the output of one phase of processing before going on to the next. But this can lead to an allocation that is duplicative of the memory allocations you must already make for the finished result, and maybe cannot be freed until the entire pipeline is done.
  - **Replace frequent, small allocations** by a longer-lived one covering the entire workflow. The result is not very idiomatically Go-like, but can have a huge impact.
  - Avoid generic libraries that come with **unpredictable memory costs**. Go’s reflection capabilities are great and let you build powerful tools. But, using them often results in costs that are not easy to pin down or control. Idioms as simple as passing in a slice rather than a fixed-sized array can have performance and memory costs. Fortunately, Go code is very easy to generate using the standard library’s go/ast and go/format packages.
* After all, my impression: We must profile, know bottleneck, read code, apply a fix, run benchmark, and loop. That is it.

## [Siloscape: The Dark Side of Kubernetes - Container Journal](https://containerjournal.com/features/siloscape-the-dark-side-of-kubernetes/)
### 2021/11/04

* siloscape: https://unit42.paloaltonetworks.com/siloscape/

## [Insecure by Default - Kubernetes Networking | Alcide](https://www.alcide.io/insecure-by-default-kubernetes-networking/)
### 2021/11/03

* Kubernetes pods by default have CAP_NET_RAW capability. It means they can open sockets and inject malcious packets into the Kubernetes network.
    * The typical threat scenario here is that an attacker has managed to take over one pod (e.g. via an application vulnerability) and wants to move laterally in the cluster to other pods.
    * Alternatively, the attacker may want to remain on the same pod but escalate their privileges to cluster-wide permissions via attacks directly against the host.
* Mitigation:
    * Drop CAP_NET_RAW. Use Kubernetes admission controller or OPA GateKeeper to prevent deploying pods with CAP_NET_RAW by developers.
    * Monitor and restrict traffic between pods by CNI or microservice firewalls.

## [How Discord Stores Billions of Messages | by Stanislav Vishnevskiy | Discord Blog](https://blog.discord.com/how-discord-stores-billions-of-messages-7fa6ec7ee4c7)
### 2021/11/03


* Originally Discord used mongoDB for the primary datastore, but as expected, many issues appered at its scale.
* They wanted to migrate to new database
* read/write patterns:
    * reads were extremely random
    * read/write ratio was about 50/50
    * Voice chat heavy Discord servers send almost no messages. This means they send a message or two every few days. In a year, this kind of server is unlikely to reach 1,000 messages. 
    * Private text chat heavy Discord servers send a decent number of messages, easily reaching between 100 thousand to 1 million messages a year. The data they are requesting is usually very recent only.
    * Large public Discord servers send a lot of messages. They have thousands of members sending thousands of messages a day and easily rack up millions of messages a year. They almost always are requesting messages sent in the last hour and they are requesting them often.
* Requirements definitions:
    * Linear scalability
    * Automatic failover
    * Low maintenance
    * Proven to work
        * not too new technology
    * Predictable performance
        * Do not want to cache messages in Redis
    * Not a blob store
    * Open source
* Cassandra was the only option which meets all the requirements
* Cassandra is KKV store; First K identifies node and location on the disk. The second is the clustering key which identifiesa row in a partition
    * A partition is something like ordered dictionary
* While migration, double write to MongoDB and Cassandra is made
* Cassandra is AP database; it is anti-pattern to read-before-write in Cassandra. What Cassandra does is essentially an upsert. You can write to any node and it will resolve conflicts automatcally using "last write wins" semantics.
    * So, in case a user edits a message at the same time as another user deletes the same message, because Cassandra's write is upsert, the row becomes empty except the primary key and text.
    * Possible solution was: 1. write the whole message back when editing the message. 2. Delete a row if a message corruption is figured out.
    * Discord chose the second option; delete messages which lacks a required column.
* Six months after Cassandra production service in, it because unresponsive
    * They noticed 10 sec GC constantly is happening.
    * The Puzzles and Dragons Subreddit public Discord serverwas the culprit. Apparently, they deleted millions of messages and only one message was left.
    * Because Cassandra does delete as soft-delete (called Tombstone), When a user loaded this channel, even though there was only 1 message, Cassandra had to effectively scan millions of message tombstones. As a result, it generates garbage faster than the JVM collect it.
    * They did 2 things; changed the tombstone lifespan from 10 days to 2 days. changed the query code to track empty buckets and avoid them in the future for a channel.

## [Practical API Design at Netflix, Part 1: Using Protobuf FieldMask | by Netflix Technology Blog | Netflix TechBlog](https://netflixtechblog.com/practical-api-design-at-netflix-part-1-using-protobuf-fieldmask-35cfdc606518)
### 2021/11/02


* We sometimes want to specify which fields are needed/not needed when calling remote API- like GraphQL.
* It is important because remote calls are not free. It imposes extra latency, error probability is increased, and consumes network bandwidth
* In JSON API, we can use sparse fieldsets
    * https://jsonapi.org/format/#fetching-sparse-fieldsets
* In gRPC, we can use FieldMask

```proto
// feld_mask.proto
message FieldMask {
  // The set of field mask paths.
  repeated string paths = 1;
}
```

* Let's say we want to call `GerProduction` API but we don't need full response from the API, we can use FieldMask like this:

```
import "google/protobuf/field_mask.proto";

message GetProductionRequest {
  string production_id = 1;
  google.protobuf.FieldMask field_mask = 2;
}
```

* If a client needs only `title` and `format`, it can build a request like this:

```
FieldMask fieldMask = FieldMask.newBuilder()
    .addPaths("title")
    .addPaths("format")
    .build();

GetProductionRequest request = GetProductionRequest.newBuilder()
    .setProductionId(LA_CASA_DE_PAPEL_PRODUCTION_ID)
    .setFieldMask(fieldMask)
    .build();
```

* Serverside implementation:

```
private static final String FIELD_SEPARATOR_REGEX = "\\.";
private static final String MAX_FIELD_NESTING = 2;
private static final String SCHEDULE_FIELD_NAME =                                // (1)
    Production.getDescriptor()
    .findFieldByNumber(Production.SCHEDULE_FIELD_NUMBER).getName();

@Override
public void getProduction(GetProductionRequest request, 
                          StreamObserver<GetProductionResponse> response) {

    FieldMask canonicalFieldMask =                                               
        FieldMaskUtil.normalize(request.getFieldMask());                         // (2) 

    boolean scheduleFieldRequested =                                             // (3)
        canonicalFieldMask.getPathsList().stream()
            .map(path -> path.split(FIELD_SEPARATOR_REGEX, MAX_FIELD_NESTING)[0])
            .anyMatch(SCHEDULE_FIELD_NAME::equals);

    if (scheduleFieldRequested) {
        ProductionSchedule schedule = 
            makeExpensiveCallToScheduleService(request.getProductionId());       // (4)
        ...
    }

    ...
}
```


## [Building Uber’s Fulfillment Platform for Planet-Scale using Google Cloud Spanner](https://eng.uber.com/building-ubers-fulfillment-platform/)
### 2021/11/02


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


## [Alerting on SLOs like Pros | SoundCloud Backstage Blog](https://developers.soundcloud.com/blog/alerting-on-slos)
### 2021/10/31

 * Because we are not Google, we don't need to follow Google way while we can learn from it
 * Originally, because a site was running on a single server, you've had to wake someone up when the server gets down
 * Nowadays, a site run on many servers and some of them must be down at any time. So it doesn't make sense to wake someone up when a single server gets down
 * Paging alerts must be strongly related to the user experience outage (must be based on symptoms, not cause). If not, it should not be a page
   * Additionally, alerts should be urgent and actionable
   * Of course, some alerts should be sent even if it doesn't affect to users, but it should be just on dashboard, not a page
* SLO should be symptoms-based alerts.
* Google tells us  in SRE book Chapter 6: “We combine heavy use of white-box monitoring with modest but critical uses of black-box monitoring.”
* In SoundCloud they mostly use white-box monitoring because:
  * “for not-yet-occurring but imminent problems, black-box monitoring is fairly useless.”
  * Because tail latency is more crucial in distributed system, even a small error can violate the SLO.
  * In microservice architecture, “one person’s symptom is another person’s cause."
* Didn't read alerts setting part.

## [Infrastructure Observability for Changing the Spend Curve - Slack Engineering](https://slack.engineering/infrastructure-observability-for-changing-the-spend-curve)
### 2021/10/30

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

## [Configuring sql.DB for Better Performance – Alex Edwards](https://www.alexedwards.net/blog/configuring-sqldb)
### 2021/10/01


Configuring sql.DB for Better Performance
https://www.alexedwards.net/blog/configuring-sqldb

* sql.DBはデータベースコネクションのプールである
    * 使用中のコネクションと、アイドルのコネクションが含まれている
* sql.DBを使う時、まずアイドルなコネクションがプールに存在するかを調べる
    * あれば、それを使う
    * なければ、コネクションを新たに張ってそれを使う
* SetMaxOpenConnsとは？
    * プールの中のコネクション数にはデフォルトでは制限がない
        * SetMaxOpenConnsはin-useとidleのコネクションの数の合計を制限する
            * なので、例えばSetMaxOpenConns(5)して、5つコネクションが全てin-useな状態で新たなコネクションが必要な場合、どれかがidleになるのを待つ。新たにコネクションを張らない。
            * 上の記事によればこう (コードはこれ https://gist.github.com/alexedwards/5d1db82e6358b5b6efcb038ca888ab07)

```
BenchmarkMaxOpenConns1-8                 500       3129633 ns/op         478 B/op         10 allocs/op
BenchmarkMaxOpenConns2-8                1000       2181641 ns/op         470 B/op         10 allocs/op
BenchmarkMaxOpenConns5-8                2000        859654 ns/op         493 B/op         10 allocs/op
BenchmarkMaxOpenConns10-8               2000        545394 ns/op         510 B/op         10 allocs/op
BenchmarkMaxOpenConnsUnlimited-8        2000        531030 ns/op         479 B/op          9 allocs/op
PASS
```

* SetMaxIdleConnsとは？
    * プール内のアイドルなコネクションの数を制限する。デフォルトでは2
    * 理論上、アイドルコネクションの数が多ければパフォーマンスは向上するはずである
    * 上の記事によればこう

```
BenchmarkMaxIdleConnsNone-8          300       4567245 ns/op       58174 B/op        625 allocs/op
BenchmarkMaxIdleConns1-8            2000        568765 ns/op        2596 B/op         32 allocs/op
BenchmarkMaxIdleConns2-8            2000        529359 ns/op         596 B/op         11 allocs/op
BenchmarkMaxIdleConns5-8            2000        506207 ns/op         451 B/op          9 allocs/op
BenchmarkMaxIdleConns10-8           2000        501639 ns/op         450 B/op          9 allocs/op
PASS
```

* MaxIdleConnを0にするといちいちコネクションを貼り直すので極めて遅い
* 1にするだけでもずいぶん変わった
* じゃあアイドルコネクションをたくさん確保すれば良いのか？
    * 場合による
    * アイドルコネクションを確保するコストはメモリを食ってしまうこと
    * アイドルコネクションは使われていないとはいえ結局データベースとつながりっぱなしではあるので
* コネクションがあまりに長い時間アイドルだったら切る、ということもできる
    * 例えばMySQLのwait_timeoutがそれを示す (デフォルトでは8h)
    * sql.DBの場合、コネクションはgracefullyにハンドルされる。コネクションがサーバーサイドから切られてもGoはそれを知るはずがないので、死んだコネクションを使おうとしたら2回リトライしたあとそれをプールから取り除き、新たにコネクションを張るような動きになる
    * なので、アイドルコネクションを持ちすぎるとより多くのリソースが使われてしまうかも。
    * 本当に使う分のアイドルコネクションだけ持つのが望ましい
* MaxIdleConnsはMaxOpenConnsより小さくなければならない (当たり前のことだが)
* SetConnMaxLifetime(d time.Duration)は、コネクションがどれだけの間使用可能かを定義する
* db.SetConnMaxLifetime(time.Hour)とすると、作られてから1時間でコネクションは 'expire' する
    * コネクションが1時間残り続けることを保証するものではない
    * コネクションは1時間経っても使われ続けていることはあり得る。1時間以上経ってから使用開始されることはないが。
    * アイドルになってから1時間ではなく、作られてから1時間
    * expireしたコネクションをプールから除外するのは1秒に1回実行される
* ConnMaxLifetimeが短いと、すぐにexpireして、より多くのコネクションが新たに張られなければならなくなる
* 上の記事によればこう

```
BenchmarkConnMaxLifetime100-8               2000        637902 ns/op        2770 B/op         34 allocs/op
BenchmarkConnMaxLifetime200-8               2000        576053 ns/op        1612 B/op         21 allocs/op
BenchmarkConnMaxLifetime500-8               2000        558297 ns/op         913 B/op         14 allocs/op
BenchmarkConnMaxLifetime1000-8              2000        543601 ns/op         740 B/op         12 allocs/op
BenchmarkConnMaxLifetimeUnlimited-8         3000        532789 ns/op         412 B/op          9 allocs/op
PASS
```

* MaxConnLifetimeを設定する場合、コネクションがどのくらい頻繁にexpire -> recreateされるのかを意識する
* 100のコネクションがあって1分のlifetimeだと、1秒に1.67のコネクションがexpireする。これがあまりに大きいとパフォーマンスが悪化しうる

http://dsas.blog.klab.org/archives/2018-02/configure-sql-db.html

ポイント

* MySQLではwait_timeoutで接続がサーバから切られた場合、MySQLドライバはそれに気づかず、クエリを終わったコネクションで送ろうとしてしまう
* SetConnMaxLifetimeを短めに設定しておけばそういうことが起こらない
* 副次的な効果は:
    * コネクションがすぐ切れるのでサーバの増減しやすくなる
    * DBのフェイルオーバーがしやすくなる
    * MySQLをオンラインで設定変更しても古い設定のコネクションが残りにくく鳴る
* SetMaxOpenConnsは必ず設定する。負荷が大きい時に新規コネクションを作らない程度の値にする
* SetMaxIdleConnsはおそらくSetMaxOpenConnsと同じにしておけば良いと思う、減らす理由がない
* SetConnMaxLifetimeは、何秒に1回再接続するか？で考える。1秒に1回程度なら殆どの場合負荷は問題にならない


Serving at localhost:8080

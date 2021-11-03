type: input
timestamp: 2021-11-03 22:23:03
url: https://blog.discord.com/how-discord-stores-billions-of-messages-7fa6ec7ee4c7
lang: en
---


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

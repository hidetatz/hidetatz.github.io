type: input
timestamp: 2021-11-24 09:45:22
url: https://dropbox.tech/application/how-dropbox-replay-keeps-everyone-in-sync
lang: en
---

* Dropbox Replay defines happend-before relationship to each message and processes them one by one
  - Using mutex, each message is received and broadcasted one by one and next message is never processed if the prev message is still in progress
* Internally the sync service must handle distributed lock? or it might be using RDB transaction for this

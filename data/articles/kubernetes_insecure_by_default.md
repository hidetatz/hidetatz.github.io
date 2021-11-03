type: input
timestamp: 2021-11-03 23:08:22
url: https://www.alcide.io/insecure-by-default-kubernetes-networking/
lang: en
---

* Kubernetes pods by default have CAP_NET_RAW capability. It means they can open sockets and inject malcious packets into the Kubernetes network.
    * The typical threat scenario here is that an attacker has managed to take over one pod (e.g. via an application vulnerability) and wants to move laterally in the cluster to other pods.
    * Alternatively, the attacker may want to remain on the same pod but escalate their privileges to cluster-wide permissions via attacks directly against the host.
* Mitigation:
    * Drop CAP_NET_RAW. Use Kubernetes admission controller or OPA GateKeeper to prevent deploying pods with CAP_NET_RAW by developers.
    * Monitor and restrict traffic between pods by CNI or microservice firewalls.

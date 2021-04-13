Report of Japan container days v18.04---2018-04-22 10:26:27

I've been to Japan container days at April 19th, 2018 in Tokyo.

[Japan container days](https://containerdays.jp/) is conference of container (almost kubernetes, k8s).  

I've been there and now want to write some impressions of sessions.  

Please note that almost all sessions were spoken in Japanese.  
I am not going to translate content of sessions into English, I will write only my own impressions.  

### Technology about private container foundation AKE by CyberAgent

original title is: サイバーエージェントにおけるプライベートコンテナ基盤AKEを支える技術

They spoke about building kubernetes cluster using OpenStack.  
They created original kubernetes as a service(kaas) on their onpremise.  
And, use it at production environment.  
I was surprised at creating own kaas, and using it at production.  
Because, if I use k8s at production, it's good to use with EKS, GKE, or other services.  

[Slide](https://speakerdeck.com/masayaaoyama/saibaezientoniokerupuraibetokontenaji-pan-akewozhi-eruji-shu) (in japanese)

### Kubernetes x Paas -- Challenge for NoOps of container applications

original title is: Kubernetes x PaaS — コンテナアプリケーションのNoOpsへの挑戦

NoOps means:

* Self healing
* In-Flight renewing
* Adaptive scale

And, k8s provides features of them. But, it make application on container statefull.  
Application on any container should be stateless.  
For it, We should use any services provided by public cloud vendor.
For example:  

* Open service broker API
* Service Catalog
* Kaas (EKS, GKE, AKS)
* Service Mesh
* Istio
* Helm

Honestly, I didn't even know about open service broker, service catalog, istio, helm.  
I want to use it with EKS or GKE and write something about this blog.

[Slide](https://www.slideshare.net/yokawasa/kubernetes-x-paas-noops) (in japanese)

### Face against microservices with Istio

original title is: Istioと共にマイクロサービスに立ち向かえ!

Microservices architecture is not a silver bullet.  
It will cause problems like this:

* difficult to manage traffic between applications
* difficult to manage keys and certificate
* difficult to glasp whole systems
* difficult to know what will be occured at failure or incident

[Istio](https://istio.io/) will help you with these annoying things.

[Slide](https://speakerdeck.com/ladicle/istiotogong-nimaikurosabisunili-tixiang-kae) (in japanese)


### Impressions

I have never operate kubernetes at production. (I only use container with AWS ECS.)  
So, I was very surprised at many services are on kubernetes cluster at production environments.  
I never know about many useful services around containers and k8s.  
I like docker, and any other container technology and I want to use them.  
Now, I want to use it and operate more simple and useful our systems.  

[Slides and entries](https://medium.com/@yukotan/japan-container-days-v18-04-%E3%81%AE%E8%B3%87%E6%96%99-4f380fb7b696)

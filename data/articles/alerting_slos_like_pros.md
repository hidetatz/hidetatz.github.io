type: input
timestamp: 2021-10-31 23:15:34
url: https://developers.soundcloud.com/blog/alerting-on-slos
lang: en
---

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

type: input
timestamp: 2021-11-04 14:59:21
url: https://www.akitasoftware.com/blog-posts/taming-gos-memory-usage-or-how-we-avoided-rewriting-our-client-in-rust
lang: en
---

* The lessons learned are interesting:
  - **Reduce fixed overhead**. Go’s garbage collection ensures that you pay for each live byte with another byte of system memory. Keeping fixed overhead low will reduce resident set size.
  - **Profile allocation, not just live data**. This reveals what is making the Go garbage collector perform work, and spikes in memory usage are usually due to increased activity at those sites.
  - **Stream, don’t buffer**. It’s a common mistake to collect the output of one phase of processing before going on to the next. But this can lead to an allocation that is duplicative of the memory allocations you must already make for the finished result, and maybe cannot be freed until the entire pipeline is done.
  - **Replace frequent, small allocations** by a longer-lived one covering the entire workflow. The result is not very idiomatically Go-like, but can have a huge impact.
  - Avoid generic libraries that come with **unpredictable memory costs**. Go’s reflection capabilities are great and let you build powerful tools. But, using them often results in costs that are not easy to pin down or control. Idioms as simple as passing in a slice rather than a fixed-sized array can have performance and memory costs. Fortunately, Go code is very easy to generate using the standard library’s go/ast and go/format packages.
* After all, my impression: We must profile, know bottleneck, read code, apply a fix, run benchmark, and loop. That is it.

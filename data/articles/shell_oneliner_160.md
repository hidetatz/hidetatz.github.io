type: input
timestamp: 2021-11-11 00:03:03
url: https://github.com/shellgei/shellgei160
lang: en
---

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

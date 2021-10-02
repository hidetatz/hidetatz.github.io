Math with fun #1: Factorization in easy way---2020-02-15 22:08:14

Let's see next formula.

$$
  6x^2 + 23x - 48
$$

If you are asked to factorize this formula, how will you do it?
Because this is quadratic equation, it will be this form.

$$
  (ax + b)(cx + d)
$$

Now, we have to know what a, b, c and d are.
This is the result of factorization of above formula.

$$
  (3x + 16)(2x - 3)
$$

There are some ways to solve this problem, but one popular way is writing this diagram.

```
3   16
  ×
2   -3
```

Multiply 3 and 2 is 6.
Multiply 16 and -3 is -48.

Now, 2 numbers arranged diagonally (3 and -3, 2 and 16) can represent 23.

$$
3 * -3 + 2 * 16 = 23
$$

To do factorization in this way, the steps are like this:

* Write 2 numbers which becomes a coefficint of $x^2$ (6) if they are multiplied
* Write 2 numbers which becomes constant term (-48) if they are multiplied
* If the result of the sum of diagonal multiplication ($3 * -3 + 2 * 16 = 23$) is the same as a coefficint of $x$ (23), then the answer is found.
  - If they are different, try another multiplication

The problem of this solution is that we are not sure how many times we have to try.
For example, divisors of 6 are 1, 2, 3, 6. So, this diagram can be this 2 patterns:

```
6

1
```

or

```
3

2
```

48 is $2^4 * 3$ . So, it has more patterns.
$(1, 48)$, $(2, 24)$, $(3, 16)$, $(4, 12)$, $(6, 8)$.

Because of above, this diagram can be 10 patterns ($2 * 5$).
Also, we have to take care of - (minus). It becomes more.
But we can simplify this method.

## Don't need to examine all the cases

To get straight to the point, we can find the answer by 2 times trial at most.

First, we need to check (6, 1) can be an answer.

```
6

1
```

Now, we don't need to check all the cases for -48. Only possible answer is (1, 48).

```
6 +-1

1 -+48
```

We can ignore other cases like (2, 24). Why?

If (6, 1) and (2, 24) works, factorized answer will be like this:

$$
(6x \pm 2)(x \mp 24)
$$

But, this can be still factorized:

$$
2(3x \pm 1)(x \mp 24)
$$

If this can be an answer, given formula must be able to be factorized by 2. However, it's obviously impossible.
In other words, left number and right number must be *coprime*. In this case, 6 and 2 is not coprime. So it never gets an answer.

```
6 +-1

1 -+48
```

This also looks not leading us to an answer. So, let's see next case (3, 2).

```
3

2
```

48 is composed of $2 * 2 * 2 * 2 * 3$. 6 is $2 * 3$.
To make them coprime, only possible answer is (16, 3).

```
3 +-16

2 -+3
```

Now, we only need to take care about sign.

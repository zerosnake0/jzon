# jzon

![](https://github.com/zerosnake0/jzon/workflows/Test/badge.svg)

## Why another jsoniter?

The code I write here is very similar to [github.com/json-iterator/go](https://github.com/json-iterator/go),
so you may ask why reinvent the wheel.

For sure that I benefit a lot from the `jsoniter` library, but i found some inconvenience for me to use it
in some condition, for example:

- the iterator methods ReadString accepts null, there is no method which accepts exactly string.
  I have to do some extra check before calling.
- some behavior is not compatible with the standard library.

On the other hand, I also want to learn how the `jsoniter` works, so there is this repo.

## What's different from jsoniter?

Here are some of the differences:

- the iterator methods accept the exact type, for example ReadString accepts only string, not null
- the behavior is almost the same as the standard library (when an error returns, the behavior may differ
  from the standard library)
- the error of the iterator is returned instead of being saved inside iterator

Some features of `jsoniter` are not implemented, and may be not implemented in the future neither.
I choose only the ones I need to implement.

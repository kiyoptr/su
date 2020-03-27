Sharder mixes data with random values and also mixes data chunks by randomly changing their offsets.

To build a sharder call `NewSharder()` with a seed and use `Add` method to add data. Use the return value of `Add`
to reference your sharded data later.
When you're done adding data, call `Shard` method to build shard result and then call `RawData` on shard result to
get the result bytes that contains mixed data.

Use `GetData` with shard result's `RawData` output and `Add`'s return value to retrieve your data.
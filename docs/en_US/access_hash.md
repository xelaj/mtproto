# Getting access_hash

Telegram is shifting issues of accessing entities to clients, and there's no way
to change that.

For some requests you need to make `InputUser`, `InputChannel`, `InputMedia` you
need a special parameter `access_hash`. The documentation does not specify how
it is organized, nor how to get it clearly. **This article describes the found
ways to get access_hash**

## Purpose of access hash

Apparently, access_hash is some kind of anti-hashing protection (so that it
would be impossible to search all channel/user/chat IDs in order and get
information about all users) and is a hash from a combination of user id, id of
the object to which the hash refers and a key that telegram servers have, but
nothing is known about the hash algorithm or what exactly is hashed. Perhaps the
answer is in [MTProxy](https://github.com/TelegramMessenger/MTProxy), but there
is no clear understanding of where to look to get the hash.

[One of answers in Stack Overflow](https://stackoverflow.com/questions/46736549/telegram-channel-how-to-get-access-hash)

## Optimization ways

One option to avoid wasting time searching for a hash is to cache it in a
database like Redis or Memcached (simple SQL tables also a good solution, just
quite slower: `access_hash` is just key-value). Please note that **DIFFERENT**
user accounts and apps have **DIFFERENT** `access_hash` values, so it is
recommended to either get public information from several _main_ accounts (from
which you collect information), and store it in your storage.

**Important:** Even that it _looks like_ collecting hashes are hard, storing,
for example 10 million hashes will consume roughly 290 megabytes, (roughly:
`int64(user_id) + int32(app_id) + int64(entity_id) + int8(entity_type) +
int64(access_hash)`, something about 29 bytes per hash value). 100 million
(nearly total user count of Teleram) will consume 2.9 gigabytes. Is it a lot?
Definitely it is, especially since `access_hash`'es contains no useful info at
all. But at the same time: imagine, how much memory you might need

## Getting hashes

> [!TIP]
> General recommendation: **cache everything you can touch.** Whole API design
> is extremely unintuitive, **so, be as much greedy as possible**.


The list of methods is most likely not final, if you have additional
information, or you have found a new method not listed in this article, please,
create a PR with a description.

### Passive collecting

Most simple way: just cache all hashes from update events you receive. It's most
simplest way: e.g. in an account of 116 supergroups we received 8600 access
hashes.

### `InputUser``

#### Resolve by nickname

contacts.resolveUsername() -> Contacts.ResolvedPeer

`Contacts.ResolvedPeer.users` contains a slice of 1 user for whom the resolution
occurred. User (user constructor) stores `id` and `access_hash`

### InputChannel

#### Resolve by invite link

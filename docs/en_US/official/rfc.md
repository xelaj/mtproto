
# RFC

> **IMPORTANT NOTE:** telegram didn't provide any specification of its mtproto spec at all. But we think that it's so important to document interfaces of any app. So this RFC actually is not **exact** rfc, it's just rework of official protocol documentation in single page.



## message format

### unencrypted data

```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                          auth_key_id                          |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                   auth_key_id (continue...)                   |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|  todo
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```


### Message Key (msg_key)

The message key is a byte sequence which contains set of bytes that allows you to check the correctness of message decoding with padding correction: for example, if the size of the encrypted message is not a multiple of the size of the aes block (128 bits, or in other words, 16 bytes), then thanks to the message key, you can decrypt the message, calculate its message key, and, if it's correct, guarantee that message decoded successfully.

#### In MTProto 2.0

To get message_key of the message, you must follow these steps:

1) get last 32 bytes of auth_key_id
2) Take all bytes of **all encrypted message**, including its headers **and** trailing padding: salt, session id, message id, seqno and padding which size can be from 12 to 1024 bytes, add it next to previous step
3) get SHA-256 hash of these bytes
4) skip 8 bytes of resulted hash and get next 16 bytes, last 8 bytes will be skipped
5) this 16 bytes (128 bits) result will be message key




### Modes


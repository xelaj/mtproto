# Notes about implementation

## Technical problems you might get stuck

* <details>
  <summary><code>-404</code> error on invalid padding size</summary>

  > [!TIP]
  >
  > For some reason, instead of making additional error (to explain developers
  > what's wrong and force them to use 12-1024 random padding). Instead,
  > Telegram server (canonical implementation of mtproto) returns -404 to EVERY
  > error that it might get while parsing message.
  >
  > **How to handle that:** write few tests and include additional check that
  > padding size is strictly between 12 and 1024 bytes, including for decrytion
  > functions.

  </details>

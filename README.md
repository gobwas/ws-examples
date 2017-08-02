# [ws](https://github.com/gobwas/ws) examples

[![website][website-image]][website-url]

> Example applications written in Go with `github.com/gobwas/ws` inside.

# Applications

- [x] [Chat](https://github.com/gobwas/ws-examples/tree/master/src/chat)
- [ ] Chat CLI
- [ ] Twitter hashtag watcher

# Notes

Currently these commands are developed:
- `bin/chat` the chat application, which is listening raw tcp socket and
  handles [jsonrpc]-like messages.
- `bin/proxy` proxy that used for two purposes. First of all, to serve static
  files for chat ui. Second and technical one is to proxy `/ws` requests to
  running chat app. This is done only for running on heroku, where only one port
  is able to be exported.

All commands can be built by `make *` or by just `make`.

Chat application deployed [here][website-url].

[website-image]: https://img.shields.io/website-up-down-green-red/http/vast-beyond-95791.herokuapp.com.svg?label=running-example
[website-url]:   https://vast-beyond-95791.herokuapp.com/#!/chat


#!/bin/sh

chat_port=3333

./bin/chat -listen=":${chat_port}" &
chat_pid=$!
for t in $(yes "1" | head -n 5); do
	ls -la /proc/${chat_pid}/fd/
	cat /proc/net/tcp
	cat /proc/net/tcp6

	if $(cat /proc/${chat_pid}/net/tcp | grep --quite LISTEN); then
		break
	fi
	sleep $t
done

./bin/proxy -listen=":${PORT}" -chat_addr=":${chat_port}"

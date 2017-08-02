#!/bin/sh

chat_port=3333

./bin/chat -listen=":${chat_port}" &
chat_pid=$!
for t in $(yes "1" | head -n 5); do
	sock=$(ls -la /proc/${chat_pid}/fd/ | fgrep 'socket' | awk '{print $4}')
	if [ ! -z "$sock" ]; then
		break
	fi
	sleep $t
done

./bin/proxy -listen=":${PORT}" -chat_addr=":${chat_port}"

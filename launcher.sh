#!/bin/bash
# Start indicated number of servers and the client
# ignore servers output

trap "exit" INT TERM		# Kill servers
trap "kill 0" exit			# after exit
declare -i port=1234

for ((i=0; i<=$1; i++))		# Start workers
do
	./worker $((port+i)) > /dev/null &		# Multiprocess
done
sleep 0.01s					# Wait for servers to start
./master $@					# Start client
#!/bin/bash
PORT=${1:-6108}

if [ "${_SERVICE_HOST_PORT}" == "" ]; then
	__IGNORE=1 #do nothing
else
	PORT=${_SERVICE_HOST_PORT}
fi

# prepare service
CAN_START_SERVICE="true"

if [ "${CAN_START_SERVICE}" == "false" ]; then
	exit 2
fi

# start service
echo '>>>>>'
echo "mode: ${_SERVICE_MODE}"
echo '<<<<<'
if [ "$_SERVICE_MODE" == "deploy" ]; then
	echo "deploy now"
	if [ "${_TRACING_MODE}" == "disable" ]; then
		echo "[jaeger] DISABLE agent"
	else
		echo "[jaeger] start agent"
		./agent-linux --collector.host-port=${_JAEGER_COLLECTOR_HOST}:14267 --processor.zipkin-compact.server-host-port 127.0.0.1:5775 &
	fi
    ./gskep
else
	echo "dev now"
	gin run
fi

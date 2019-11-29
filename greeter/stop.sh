ps -ef | grep greeter | grep -v grep | awk '{print $2}' | xargs kill -9

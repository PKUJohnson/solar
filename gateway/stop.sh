ps -ef | grep gateway | grep -v grep | awk '{print $2}' | xargs kill -9

ps -ef | grep user | grep -v grep | awk '{print $2}' | xargs kill -9

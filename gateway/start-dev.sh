export LD_LIBRARY_PATH=/usr/local/lib/:$LD_LIBRARY_PATH
export CONFIGOR_ENV=dev
nohup ./gateway >/dev/null 2>&1 &


export STAGE="local"

export LOCAL_STORE_PATH="/store"
export LOCAL_HTTP_URL=":18061"

cd ./bin_dev
while true; do
    ./wartank_dev
    sleep 1
done

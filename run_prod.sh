export STAGE="local"

export LOCAL_STORE_PATH="/store"
export LOCAL_HTTP_URL=":18060"

cd ./bin
while true; do
    ./wartank_prod
    sleep 1
done

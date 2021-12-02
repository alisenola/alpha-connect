docker build --network=host -f alpha-connect.Dockerfile . -t alpha-connect
if [ ! -z "$1" ]
then
   docker tag alpha-connect:latest gcr.io/patrick-206008/alpha-connect:$1
   docker push gcr.io/patrick-206008/alpha-connect:$1
fi

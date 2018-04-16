sh ./generateData.sh 5000 > ./records.sql
docker run -ti -v ${PWD}:/root/ --link aerospike:aerospike aerospike/aerospike-tools aql -h aerospike --no-config-file --file /root/records.sql
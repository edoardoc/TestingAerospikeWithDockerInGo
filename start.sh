docker-compose down --rmi all
docker-compose build
docker-compose up -d

sleep 5
sh populate.sh
rm cibucks_data/*
docker-compose down --rmi all
docker-compose build
docker-compose up -d

sleep 5
sh populate.sh
sleep 1
docker logs -f golangserver

./network.sh down
docker stop logspout 
docker rm logspout 
docker volume rm $(docker volume ls)


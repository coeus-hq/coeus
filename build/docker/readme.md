**Stop/Remove the container if running**  
docker stop coeus  
docker rm coeus  

**Remove the image**  
docker rmi -f `docker images --format="{{.Repository}} {{.ID}}" | grep "^coeus " | cut -d' ' -f2`  
 
(docker system prune --force --all)

**Build new image and verify**  
docker build --progress=plain -t coeus  -f coeus.dockerfile .  
docker images  

**Run container and verify**  
docker run --name coeus -h coeus-docker -d  -p 5432:5432 -p 8080:8080  -t coeus  
docker ps  
docker logs coeus  

**Connect with psql**  
export PGPASSWORD=coeus  
psql -h localhost -U coeus -d coeus  

**Optionally attach with the running container**  
docker exec -it coeus  /bin/bash  

Open http://localhost:8080/  


docker tag <image id> asim95/coeus:latest
  357  docker push asim95/coeus



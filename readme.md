# Social tournament service
used env - Ubuntu 16.04  
sudo apt-get update  
sudo apt-get upgrade -y  
sudo apt-get install git -y  
sudo apt-get install docker.io -y

sudo docker pull mongo
sudo docker run --name mongo-cont -d mongo
sudo docker exec -it mongo-cont mongo

sudo git clone https://github.com/OlexandrPodustov/sts
cd sts  
sudo docker build -t podustov .
sudo docker run --rm -it -p 8080:8080 podustov cmd


http://localhost:8080/fund?playerId=P1&points=300 
http://localhost:8080/take?playerId=P1&points=300
http://localhost:8080/announceTournament?tournamentId=1&deposit=1000
http://localhost:8080/joinTournament?tournamentId=1&playerId=P1&backerId=P2&backerId=P3
POST http://localhost:8080/resultTournament
http://localhost:8080/balance?playerId=P1
http://localhost:8080/reset

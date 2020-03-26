# social tournament service
docker pull mongo  
docker run --name mongo-cont -d mongo  

git clone git@github.com:OlexandrPodustov/sts.git  
cd sts  
docker build -t podustov .  
docker run --rm -it -p 8080:8080 podustov cmd


GET http://localhost:8080/fund?playerId=P1&points=300  
GET http://localhost:8080/take?playerId=P1&points=300  
GET http://localhost:8080/announceTournament?tournamentId=1&deposit=1000  
GET http://localhost:8080/joinTournament?tournamentId=1&playerId=P1&backerId=P2&backerId=P3  
POST http://localhost:8080/resultTournament  
GET http://localhost:8080/balance?playerId=P1  
GET http://localhost:8080/reset
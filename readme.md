# Social tournament service
used env - Ubuntu 16.04  
sudo apt-get update  
sudo apt-get upgrade -y  
sudo apt-get install git -y  
sudo apt-get install docker.io -y

git clone https://github.com/OlexandrPodustov/sts  
cd sts  
docker build -t podustov .  
docker run --rm -it -p 8080:8080 podustov ./cmd

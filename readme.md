# Social tournament service
used env - Ubuntu 16.04  
sudo apt-get update  
sudo apt-get upgrade -y  
sudo apt-get install git -y  
sudo apt-get install docker.io -y

sudo docker pull mongo


sudo git clone https://github.com/OlexandrPodustov/sts
cd sts  
sudo docker build -t podustov .
sudo docker run --rm -it -p 8080:8080 podustov cmd

http://localhost:8080/reset

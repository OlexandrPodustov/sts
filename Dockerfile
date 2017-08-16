FROM ubuntu:xenial
RUN   apt-get update\
&& apt-get -y upgrade\
&& apt-get install curl -y\
&& curl -O https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz \
&& tar -C /usr/local -xzf go1.8.3.linux-amd64.tar.gz\
&& export PATH=$PATH:/usr/local/go/bin\
&& apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 0C49F3730359A14518585931BC711F9BA15703C6\
&& echo "deb [ arch=amd64,arm64 ] http://repo.mongodb.org/apt/ubuntu xenial/mongodb-org/3.4 multiverse" | tee /etc/apt/sources.list.d/mongodb-org-3.4.list\
&& apt-get update\
&& apt-get install -y mongodb-org\
&&  mkdir -p /data/db\
&& mongod -d\



EXPOSE 80
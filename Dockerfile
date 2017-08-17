FROM ubuntu:xenial
RUN   apt-get update \
&& apt-get upgrade -y \
&& apt-get install curl -y \
&& curl -O https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz \
&& tar -C /usr/local -xzf go1.8.3.linux-amd64.tar.gz


#&& apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 0C49F3730359A14518585931BC711F9BA15703C6 \
#&& echo "deb [ arch=amd64,arm64 ] http://repo.mongodb.org/apt/ubuntu xenial/mongodb-org/3.4 multiverse" | tee /etc/apt/sources.list.d/mongodb-org-3.4.list \
#&& apt-get update \
#&& apt-get install mongodb-org -y \
#&& mkdir -p /data/db \
#&& mongod &
ENV PATH=$PATH:/usr/local/go/bin

COPY . /usr/local/go/src/sts
RUN cd /usr/local/go
RUN env GOOS=linux GOARCH=amd64 go build -v sts/cmd
#RUN /usr/local/go/src/sts/cmd.exe &
RUN /usr/local/go/src/sts/cmd.exe &

EXPOSE 8080
CMD ["cmd", "/go/src/sts"]
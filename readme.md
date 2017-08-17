# Social tournament service
git clone https://github.com/OlexandrPodustov/sts
cd sts
docker build -t podustov .
docker run --rm -it -p 8080:8080 podustov ./cmd
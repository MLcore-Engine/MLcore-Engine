### infra:
#### registry:  
docker run -d -p 5005:5000 --restart=always --name registry -v /data/imageregistory/:/var/lib/registry registry:latest

#### minio:
 docker run   -p 10000:9000   -p 10001:9090   --name minio   -d   --restart=always   -e "MINIO_ACCESS_KEY=admin"   -e "MINIO_SECRET_KEY=admin12345"   -v /data/miniodata/data:/data   -v /data/miniodata/config:/root/.minio   minio/minio server   /data --console-address ":9090" -address ":9000"
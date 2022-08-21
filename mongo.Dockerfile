FROM mongo:latest
COPY resources/fixtures.mongo.js /docker-entrypoint-initdb.d/

CMD ["--wiredTigerCacheSizeGB", "0.25", "--bind_ip_all"]
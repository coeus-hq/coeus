FROM postgres:alpine

# Expose ports for coeus web and and postgres db
EXPOSE 8080
#EXPOSE 5432

COPY coeus-bin /root
COPY start.sh /root
WORKDIR /root
CMD ["/root/start.sh"]

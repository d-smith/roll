FROM voxxit/base:alpine

ENV http_proxy http://***REMOVED***:8000
ENV https_proxy http://***REMOVED***:8000

ADD vconfig /var/vconfig

RUN  wget http://dl.bintray.com/mitchellh/vault/vault_0.3.1_linux_amd64.zip \
  && unzip vault_0.3.1_linux_amd64.zip \
  && mv vault /usr/local/bin/ \
  && rm -f vault_0.3.1_linux_amd64.zip

EXPOSE 8200

ENTRYPOINT [ "vault" ]
CMD [ "server", "-config=/var/vconfig" ]

#docker run --cap-add IPC_LOCK -p 8200:8200 roll/vault
FROM voxxit/base:alpine

ENV http_proxy http://proxyhost:proxyport
ENV https_proxy http://proxyhost:proxyport

ADD vconfig /var/vconfig

RUN  wget http://dl.bintray.com/mitchellh/vault/vault_0.3.1_linux_amd64.zip \
  && unzip vault_0.3.1_linux_amd64.zip \
  && mv vault /usr/local/bin/ \
  && rm -f vault_0.3.1_linux_amd64.zip

EXPOSE 8200

ENTRYPOINT [ "vault" ]
CMD [ "server", "-config=/var/vconfig" ]

#docker run -d --cap-add IPC_LOCK -p 8200:8200 roll/vault
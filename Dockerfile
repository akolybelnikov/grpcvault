FROM scratch
MAINTAINER Andrey Kolybelnikov a.kolybelnikov@gmail.com
ADD vaultd vaultd
EXPOSE 8080 8081
ENTRYPOINT ["/vaultd"]
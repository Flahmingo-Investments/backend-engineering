FROM golang:1.16-stretch AS base

# Install doppler tool
RUN (curl -Ls --tlsv1.2 --proto "=https" --retry 3 https://cli.doppler.com/install.sh || wget -t 3 -qO- https://cli.doppler.com/install.sh) | sh

# Install system tools
RUN apt update && apt install -y unzip

# Install proto tools
ARG PB_REL=https://github.com/protocolbuffers/protobuf/releases
RUN curl -LO $PB_REL/download/v3.15.8/protoc-3.15.8-linux-x86_64.zip
RUN unzip protoc-3.15.8-linux-x86_64.zip -d /go
RUN rm protoc-3.15.8-linux-x86_64.zip
RUN chmod +x /go/bin/protoc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
RUN go get github.com/pressly/goose/v3/cmd/goose


FROM base AS local

# Argumetns
ARG USERNAME=infinito
ARG USER_UID=1000
ARG USER_GID=$USER_UID

# Install sytem tools for local target
RUN apt update && apt install -y sudo postgresql-client-common postgresql-client libpq-dev locales
RUN curl -LO https://github.com/fullstorydev/grpcurl/releases/download/v1.8.5/grpcurl_1.8.5_linux_x86_64.tar.gz
RUN tar -xvf grpcurl_1.8.5_linux_x86_64.tar.gz
RUN mv grpcurl /go/bin/
RUN rm grpcurl_1.8.5_linux_x86_64.tar.gz LICENSE

# Add User
RUN groupadd --gid $USER_GID $USERNAME
RUN useradd --uid $USER_UID --gid $USER_GID -m $USERNAME
RUN mkdir -p /home/$USERNAME/.vscode-server/extensions /home/$USERNAME/.vscode-server-insiders/extensions
RUN chown -R $USERNAME /home/$USERNAME/.vscode-server /home/$USERNAME/.vscode-server-insiders
RUN echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME
RUN chmod 0440 /etc/sudoers.d/$USERNAME

# Preserving bash history
RUN SNIPPET="source /workspace/scripts/bash" \
    && mkdir /history \
    && touch /history/.bash_history \
    && chown -R $USERNAME /history \
    && echo $SNIPPET >> "/home/$USERNAME/.bashrc"

# Install go tools
RUN go get -v github.com/ramya-rao-a/go-outline
RUN go get -v golang.org/x/tools/gopls
RUN go get -v github.com/uudashr/gopkgs/v2/cmd/gopkgs
RUN go get -v github.com/go-delve/delve/cmd/dlv
RUN go get -v github.com/go-delve/delve/cmd/dlv@master
RUN go get -v honnef.co/go/tools/cmd/staticcheck

# Fix permissions
RUN chmod -R 777 /go

# Using User
USER $USERNAME

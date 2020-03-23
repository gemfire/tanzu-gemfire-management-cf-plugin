FROM ubuntu:18.04

RUN apt-get update

# Install basic utilities
RUN apt-get install -y \
    wget \
    gnupg \
    curl \
    jq

# Install cf
RUN wget -q -O - https://packages.cloudfoundry.org/debian/cli.cloudfoundry.org.key | apt-key add -
RUN echo "deb https://packages.cloudfoundry.org/debian stable main" > /etc/apt/sources.list.d/cloudfoundry-cli.list
RUN apt-get update
RUN apt-get install -y cf-cli

# Install om
RUN curl -L "https://github.com/pivotal-cf/om/releases/download/3.0.0/om-linux-3.0.0" -o "om-linux"
RUN install om-linux /usr/bin/om

# Install bosh
RUN apt-get install -y \
    build-essential \
    zlibc \
    zlib1g-dev \
    ruby \
    ruby-dev \
    openssl \
    libxslt1-dev \
    libxml2-dev \
    libssl-dev \
    libreadline7 \
    libreadline-dev \
    libyaml-dev \
    libsqlite3-dev \
    sqlite3 \
    openssh-client \
    netcat-openbsd
RUN curl -L "https://github.com/cloudfoundry/bosh-cli/releases/download/v5.5.1/bosh-cli-5.5.1-linux-amd64" -o "bosh"
RUN chmod +x bosh
RUN install bosh /usr/bin/bosh


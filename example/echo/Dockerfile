FROM debian:jessie

RUN sed -i 's/deb.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list

#RUN apt-get update \
# && apt-get install -y locales \
# && dpkg-reconfigure -f noninteractive locales \
# && locale-gen C.UTF-8 \
# && /usr/sbin/update-locale LANG=C.UTF-8 \
# && echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen \
# && locale-gen \
# && apt-get clean \
# && rm -rf /var/lib/apt/lists/*

# Users with other locales should set this in their derivative image
#ENV LANG en_US.UTF-8
#ENV LANGUAGE en_US:en
#ENV LC_ALL en_US.UTF-8

COPY echo /usr/local/bin

CMD ["echo"]

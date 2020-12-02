FROM registry.access.redhat.com/ubi8/ubi-minimal

LABEL author="gatblau"
LABEL maintainer="onix@gatblau.org"
LABEL builder-for=PYTHON-36
LABEL tools="python"

ARG APP_NAME=artie
ARG UNAME=builder

ENV UID=1000
ENV GID=1000
ENV APP=$APP_NAME
ENV ARTIE_HOME /usr/bin

RUN microdnf update --disablerepo=* --enablerepo=ubi-8-appstream --enablerepo=ubi-8-baseos -y && \
    microdnf install shadow-utils.x86_64 && \
    microdnf install python3 && \
    groupadd -g $GID -o $UNAME && useradd -M -u $UID -g $GID $UNAME && \
    rm -rf /var/cache/yum && \
    microdnf clean all

COPY ./images/bin/output/$APP_NAME ./images/start.sh $ARTIE_HOME/

RUN ln -s $ARTIE_HOME/$APP_NAME $ARTIE_HOME/artie && \
    chmod ug+x $ARTIE_HOME/$APP_NAME && \
    chmod ug+x $ARTIE_HOME/artie && \
    chmod ug+x $ARTIE_HOME/start.sh && \
    # folder below to mount keys in container
    mkdir -p /keys

USER $UNAME

ENTRYPOINT ["sh", "start.sh"]



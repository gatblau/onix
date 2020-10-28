FROM centos:8

RUN dnf -y install unzip && \
    dnf -y install wget && \
    dnf install -y java-11-openjdk-devel

ENV MAVEN_VERSION 3.6.3
ENV JAVA_HOME /usr/lib/jvm/jre-openjdk
ENV MAVEN_HOME /usr/lib/mvn
ENV PATH $MAVEN_HOME/bin:$PATH

# download and install maven in the build container
RUN wget http://archive.apache.org/dist/maven/maven-3/$MAVEN_VERSION/binaries/apache-maven-$MAVEN_VERSION-bin.tar.gz && \
  tar -zxvf apache-maven-$MAVEN_VERSION-bin.tar.gz && \
  rm apache-maven-$MAVEN_VERSION-bin.tar.gz && \
  mv apache-maven-$MAVEN_VERSION /usr/lib/mvn

COPY artie .

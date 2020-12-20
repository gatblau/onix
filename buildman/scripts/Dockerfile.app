FROM registry.access.redhat.com/ubi8/ubi-minimal
# need to add something in order to update the CREATED property in the image manifest
RUN echo "this is the app image" >> app
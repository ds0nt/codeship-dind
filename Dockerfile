FROM progrium/busybox

ADD codeship-dind ./codeship-dind
CMD codeship-dind

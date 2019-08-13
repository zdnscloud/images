#!/bin/bash

BRANCH=$1

if [[ -z $BRANCH ]]
then
cat <<EOF

<------------------------------------------------>
    Usage: ./gen_image.sh {branch}
<------------------------------------------------>

EOF
    exit 1
else
    cat <<EOF

<------------------------------------------------>
    Pull newest image of {branch}
<------------------------------------------------>

EOF
    docker pull zdnscloud/singlecloud:${BRANCH}
    docker pull zdnscloud/singlecloud-ui:${BRANCH}
fi

cat <<EOF

<------------------------------------------------>
    Building ...
<------------------------------------------------>

EOF

cat <<'EOF' | docker build -f - -t zdnscloud/singlecloud:build-${BRANCH} --build-arg branch=${BRANCH} .
ARG branch

FROM zdnscloud/singlecloud:$branch as go
FROM zdnscloud/singlecloud-ui:$branch as js

FROM alpine:latest

LABEL zcloud/branch=$branch

RUN apk --no-cache add ca-certificates
COPY --from=go /usr/local/bin/singlecloud /usr/local/bin
COPY --from=js /www /www

EXPOSE 80
CMD ["-listen", ":80"]

ENTRYPOINT ["/usr/local/bin/singlecloud"]     
EOF

if [[ $? -eq 0 ]]
then
cat <<EOF

<------------------------------------------------>
  Image build complete.
  Build: zdnscloud/singlecloud:build-${BRANCH}
  Run: docker run --rm -p 8080:80 zdnscloud/singlecloud:build-${BRANCH}
<------------------------------------------------>

EOF
else
cat <<EOF

<------------------------------------------------>
  Image build failure.
<------------------------------------------------>

EOF
fi

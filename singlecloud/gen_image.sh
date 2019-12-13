#!/bin/bash

set -e

BRANCH=$1
UIBRANCH=${2:-$1}

if [[ -z $BRANCH ]]
then
cat <<EOF

<------------------------------------------------>
    Usage:
        ./gen_image.sh {branch}
        ./gen_image.sh {singlecloud branch} {singlecloud-ui branch}
<------------------------------------------------>

EOF
    exit 1
else
    cat <<EOF

<------------------------------------------------>
    Pull newest image of singlecloud - ${BRANCH}
    Pull newest image of singlecloud-ui - ${UIBRANCH}
<------------------------------------------------>

EOF
    docker pull zdnscloud/singlecloud:${BRANCH}
    docker pull zdnscloud/singlecloud-ui:${UIBRANCH}
fi

cat <<EOF

<------------------------------------------------>
    Building ...
<------------------------------------------------>

EOF

cat <<'EOF' | docker build -f - -t zdnscloud/singlecloud:${BRANCH}-${UIBRANCH} --build-arg branch=${BRANCH} --build-arg uibranch=${UIBRANCH} .
ARG branch
ARG uibranch

FROM zdnscloud/singlecloud:$branch as go
FROM zdnscloud/singlecloud-ui:$uibranch as js

FROM scratch

LABEL zcloud/branch=$branch
LABEL ui.zcloud/branch=$uibranch

COPY --from=go /singlecloud /
COPY --from=js /www /www

EXPOSE 80
CMD ["-listen", ":80"]

ENTRYPOINT ["/singlecloud"]
EOF

if [[ $? -eq 0 ]]
then
docker image prune -f
cat <<EOF

<------------------------------------------------>
  Image build complete.
  Build: zdnscloud/singlecloud:${BRANCH}-${UIBRANCH}
  Run: docker run --rm -p 8080:80 zdnscloud/singlecloud:${BRANCH}-${UIBRANCH}
<------------------------------------------------>

EOF
else
cat <<EOF

<------------------------------------------------>
  Image build failure.
<------------------------------------------------>

EOF
fi

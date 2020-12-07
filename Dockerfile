FROM golang:1.15.6 as builder

WORKDIR /go/src/github.com/swade1987/perimener

COPY . .

RUN make perimener-linux

FROM scratch
MAINTAINER Steve Wade <steven@stevenwade.co.uk>

ARG git_repository="Unknown"
ARG git_commit="Unknown"
ARG git_branch="Unknown"
ARG built_on="Unknown"

LABEL git.repository=$git_repository
LABEL git.commit=$git_commit
LABEL git.branch=$git_branch
LABEL build.on=$built_on

COPY --from=builder /go/src/github.com/swade1987/perimener/bin/linux/perimener .

CMD [ "/perimener" ]

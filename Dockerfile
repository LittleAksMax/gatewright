FROM --platform=$BUILDPLATFORM golang:1.27-bookworm as builder

ARG TARGETARCH

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ cmd/
COPY main.go main.go

COPY Makefile Makefile
RUN make build TARGETARCH=${TARGETARCH}


FROM python:3.14-slim-bookworm AS spec-tests

RUN apt-get update \
    && apt-get install -y --no-install-recommends sudo \
    && rm -rf /var/lib/apt/lists/* \
    && useradd --create-home --shell /bin/bash tester \
    && echo "tester ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/tester \
    && chmod 0440 /etc/sudoers.d/tester

COPY --from=builder /app/gwt /usr/local/bin/gwt
COPY spec_tests/requirements.txt /tmp/requirements.txt
RUN pip install --no-cache-dir -r /tmp/requirements.txt

COPY spec_tests/ /spec_tests/

USER tester
WORKDIR /spec_tests

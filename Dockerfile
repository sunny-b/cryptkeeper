FROM golang:latest

RUN apt-get update && \
    apt-get install -y bash zsh fish xclip direnv xvfb

# Create an entrypoint script
RUN echo '#!/bin/bash\nXvfb :1 &\nexport DISPLAY=:1\nexec "$@"' > /entrypoint.sh && \
    chmod +x /entrypoint.sh

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2

WORKDIR /workspace

COPY . .

RUN make build && mv /workspace/bin/cryptkeeper /usr/local/bin/cryptkeeper

# Set the entrypoint
ENTRYPOINT ["/entrypoint.sh"]

CMD ["make", "test"]

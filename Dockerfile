FROM golang:1.22.4-alpine

# Setup Work Directory
WORKDIR /app

# Copy Go modules and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o golangmanpower .

# Expose the port your app will run on
EXPOSE 5250

ENV PORT=${PORT}
ENV DB_USER=${DB_USER}
ENV DB_PASS=${DB_PASS}
ENV DB_NAME=${DB_NAME}
ENV DB_SERVER=${DB_SERVER}
ENV LDAP_IP=${LDAP_IP}
ENV LDAP_DNS=${LDAP_DNS}
ENV SECRET_KEY=${SECRET_KEY}

# Run the application
CMD ["./golangmanpower"]
# Use a Node.js image
FROM node:20-alpine

# Install protobuf compiler (protoc) and grpc-web plugin
RUN apk add --no-cache protobuf-dev && \
    npm install -g protoc-gen-grpc-web && \
    # Create the google/protobuf directory expected by our proto definitions
    mkdir -p /usr/include/google/protobuf && \
    # Copy google/protobuf/empty.proto if it's external, or rely on protobuf-dev package
    # For simplicity, we'll assume protobuf-dev provides common ones like empty.proto
    # If not, you might need to manually copy it in here
    cp -r /usr/local/include/google/protobuf/empty.proto /usr/include/google/protobuf/empty.proto || true


# Set the working directory inside the container
WORKDIR /app

# Copy package.json and package-lock.json (if exists)
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the project files
COPY . .

# Create the output directory
RUN mkdir -p frontend/proto

# Define a default command (you'll override this for specific tasks)
CMD ["node"]
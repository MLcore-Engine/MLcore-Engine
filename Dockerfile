From gaoxin2020/aiplatform:base


WORKDIR /app

# Copy the React project into the container
COPY ./mlcore-engine ./

# Install dependencies and build the React app
RUN cd web
RUN npm run build
RUN cd ..
RUN go mod tidy && go build -o /app/main

# Expose the port Go server will run on
EXPOSE 3000

# Run the Go app
CMD ["/app/main"]


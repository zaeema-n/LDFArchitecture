# Use Ubuntu 22.04 as the base image
FROM ballerina/ballerina:2201.11.0

# Copy the Ballerina source code from the current directory
COPY hello_world.bal .

# Build the Ballerina program
RUN bal build hello_world.bal

# Run the program
CMD ["bal", "run", "hello_world.bal"]

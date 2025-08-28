## Data Stream Service

The Data Stream Service is part of the Sensor Pool microservice system. This service acts as a sensor and sends messages periodically. This service will send message to mqtt broker.

### Setup

1. Get api key from data pool service.
2. Insert api key, sensor identity, and broker information into the .env file
   following the .env.example
3. Build the docker image, run command below:

   ```
   docker build -t <image name> .
   ```

   The command above will create an image of the app.

4. Create a container from the image and run command below:

   ```
   docker run -d --name <container name> -p 8081:8080 -v $(pwd)/.env:/build/.env <image name>
   ```

   The command above tell docker to create and run the image, in the command above we use dynamic .env for easier and more secure .env file configuration.

### REST API

This service exposes a REST API to change the frequency at which messages are sent to the MQTT broker, the default frequency is `5` second. Below is the endpoint.

```
/sensor?d=<new frequency in second>
```

The endpoint requires a parameter d, which sets the new message frequency in seconds.

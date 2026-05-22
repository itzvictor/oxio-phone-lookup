# oxio-phone-lookup

## Prerequisites

- [Go 1.26+](https://go.dev/dl/)

## Start the server locally

```bash
go run main.go
```

The server listens on port `8080`.

## Test the API locally

Look up a phone number with a country code in the number:

```bash
curl "http://localhost:8080/v1/phone-numbers?phoneNumber=%2B12125690123"
```

Look up a local number by providing a country code separately:

```bash
curl "http://localhost:8080/v1/phone-numbers?phoneNumber=2125690123&countryCode=US"
```

## Run the tests

```bash
go test ./...
```

---

## FAQ

### Technology choices

- **chi** - I am most familiar with this but any router library should be fine. 
- **oapi-codegen** - This generate Go structs from the api_spec. It creates a tight contract between the client and server.
- **nyaruka/phonenumbers** - It looks the most verbose library and it is well-maintained Go port of Google's `libphonenumber`. 
- **mockery** - Generates type-safe mocks from interfaces, used to isolate the controller layer from the service in unit tests.

### Deploying to production


2. **Containerize** - Build it in a minimal Docker image like distroless using Dockerfile
3. **Deploy** - This is lightweight and stateless, I'd just use ECS.
5. **Observability** - Add logging to CloudWatch, a health check endpoint, and metrics (e.g. Prometheus)

### Assumptions

- Area code extraction assumes 7 digits for the local number, which is accurate for NANP (North American) numbers.

### Improvements

- **Authenication** - Ideally the user should be authenticated via JWT bearer token in the Authorization header. This would make a request to a auth service/library.
- **Errors** - I didn't have time to add proper errors. The service layer should have its own error codes for each type of error and those errors are propagated up to the controller caller. For example 1000 -> 400, 1001 -> 500. Currently it always return 400 on service error but it might not be true in the future
- **E2E testing** - I would generate a Go client with oapi to mimic how a consumer (UI or another service) would interact with my service. 
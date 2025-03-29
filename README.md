## Running with Docker

```bash
docker-compose up -d
```

This will start:
- PostgreSQL database (port 5432)
- MinIO object storage (ports 9000, 9001)
- The main application (port 8080)

##  Documentation

Once the application is running, you can access the Swagger documentation at:

```
http://localhost:8080/swagger/index.html
```

The Swagger UI provides a complete API reference with endpoint details and allows you to:
- View all available API endpoints
- Test API endpoints directly from the browser
- See request/response schemas and examples
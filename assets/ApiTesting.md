# API Testing

Below are some core API endpoints for the project. You can test these APIs using cURL in the terminal or import them into Postman (it will automatically convert to a request).

## Register Account
```bash
curl -X POST http://localhost:8080/api/users/register \
-H "Content-Type: application/json" \
-d '{
  "username": "username",
  "password": "password",
  "email": "email@example.com",
  "firstName": "Capi",
  "lastName": "Giba"
}'
```

## Login and reveice a token
Note: Currently, the project supports login via email. Username support will be added once the product is stable.
```
curl -X POST http://localhost:8080/api/users/login \
-H "Content-Type: application/json" \
-d '{
  "email": "john@example.com",
  "password": "password123"
}'
```

## Get Presigned URL
```
curl -X POST http://localhost:8080/api/videos/generate-presigned-url \
-H "Content-Type: application/json" \
-H "Authorization: Bearer <YOUR_JWT_TOKEN>" \
-d '{
  "user_id": ?,
  "file_name": "my_video.mp4",
  "contentType": "video/mp4",
  "title": "Sample Video Title",
  "duration": 120
}'
```

## Others: 
```
In Progress...
```
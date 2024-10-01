# API Documentation for User Features

## 1. User Registration
- **API Endpoint**: `POST /users/register`
- **Description**: Registers a new user in the system.
- **Input** (JSON body):
    ```json
    {
        "first_name": "John",
        "last_name": "Doe",
        "username": "johndoe123",
        "email": "johndoe@example.com",
        "password": "securePassword123"
    }
    ```
- **Response**:
    - `201 Created`: User registered successfully.
    - `400 Bad Request`: Validation error.
    - `500 Internal Server Error`: Server-side issue.

## 2. User Login
- **API Endpoint**: `POST /users/login`
- **Description**: Logs in a user by validating their credentials.
- **Input** (JSON body):
    ```json
    {
        "email": "johndoe@example.com",
        "password": "securePassword123"
    }
    ```
- **Response**:
    - `200 OK`: Returns an authentication token.
    - `400 Bad Request`: Validation error.
    - `401 Unauthorized`: Invalid credentials.

## 3. Get User Details
- **API Endpoint**: `GET /users/{user_id}`
- **Description**: Retrieves user information by user ID.
- **Input** (Path parameter):
    - `user_id` (int): ID of the user to be fetched.
- **Response**:
    ```json
    {
        "id": 1,
        "first_name": "John",
        "last_name": "Doe",
        "username": "johndoe123",
        "email": "johndoe@example.com",
        "status": 1,
        "premium": true,
        "role": "User",
        "avatar": "avatar.jpg",
        "avatar_folder": "avatars/123/",
        "created_at": "2023-09-01T12:34:56Z",
        "updated_at": "2023-09-01T12:34:56Z"
    }
    ```

## 4. Update User Information
- **API Endpoint**: `PUT /users/{user_id}`
- **Description**: Updates user information (excluding avatar).
- **Input** (Path parameter & JSON body):
    - `user_id` (int): ID of the user.
    ```json
    {
        "first_name": "John",
        "last_name": "Doe",
        "username": "johndoe123",
        "email": "johndoe@example.com",
        "premium": true,
        "role": "Admin"
    }
    ```
- **Response**:
    - `200 OK`: User updated successfully.
    - `400 Bad Request`: Validation error.
    - `500 Internal Server Error`: Server-side error.

## 5. Change Password
- **API Endpoint**: `PUT /users/{user_id}/change-password`
- **Description**: Allows the user to change their password.
- **Input** (Path parameter & JSON body):
    - `user_id` (int): ID of the user.
    ```json
    {
        "old_password": "oldPassword123",
        "new_password": "newPassword123"
    }
    ```
- **Response**:
    - `200 OK`: Password changed successfully.
    - `400 Bad Request`: Validation error.
    - `500 Internal Server Error`: Server-side error.

## 6. Delete User (Soft Delete)
- **API Endpoint**: `DELETE /users/{user_id}`
- **Description**: Soft deletes a user by updating their status.
- **Input** (Path parameter):
    - `user_id` (int): ID of the user to be deleted.
- **Response**:
    - `200 OK`: User deleted successfully.
    - `500 Internal Server Error`: Server-side error.

## 7. Update Avatar
- **API Endpoint**: `PUT /users/{user_id}/update-avatar`
- **Description**: Generates a presigned URL for uploading a new avatar for the user.
- **Input** (Path parameter & query parameters):
    - `user_id` (int): ID of the user.
    - `file_name` (string): The file name for the avatar.
    - `folder` (string): The folder to store the avatar.
- **Response** (Example JSON response):
    ```json
    {
        "avatar_upload_url": "https://s3.amazonaws.com/examplebucket/avatars/123/avatar.jpg?presigned-url"
    }
    ```
    - `500 Internal Server Error`: Server-side error.

## 8. Get Avatar Download URL
- **API Endpoint**: `GET /users/{user_id}/avatar-download-url`
- **Description**: Generates a presigned URL to download the user's avatar from S3.
- **Input** (Path parameter):
    - `user_id` (int): ID of the user.
- **Response** (Example JSON response):
    ```json
    {
        "avatar_download_url": "https://s3.amazonaws.com/examplebucket/avatars/123/avatar.jpg?presigned-url"
    }
    ```
    - `500 Internal Server Error`: Server-side error.

## 9. Load Avatar
- **API Endpoint**: `GET /users/{user_id}/avatar`
- **Description**: Redirects to the avatar image using a presigned URL.
- **Input** (Path parameter):
    - `user_id` (int): ID of the user.
- **Response**:
    - `307 Temporary Redirect`: Redirects to the avatar URL.
    - `500 Internal Server Error`: Server-side error.

## 10. Get All Users
- **API Endpoint**: `GET /users`
- **Description**: Retrieves a list of all users.
- **Input**: None.
- **Response** (Example JSON response):
    ```json
    [
        {
            "id": 1,
            "first_name": "John",
            "last_name": "Doe",
            "username": "johndoe123",
            "email": "johndoe@example.com",
            "status": 1,
            "premium": true,
            "role": "User",
            "avatar": "avatar.jpg",
            "avatar_folder": "avatars/123/",
            "created_at": "2023-09-01T12:34:56Z",
            "updated_at": "2023-09-01T12:34:56Z"
        },
        {
            "id": 2,
            "first_name": "Jane",
            "last_name": "Smith",
            "username": "janesmith456",
            "email": "janesmith@example.com",
            "status": 1,
            "premium": false,
            "role": "Admin",
            "avatar": "avatar2.jpg",
            "avatar_folder": "avatars/456/",
            "created_at": "2023-09-01T12:34:56Z",
            "updated_at": "2023-09-01T12:34:56Z"
        }
    ]
    ```
    - `500 Internal Server Error`: Server-side error.

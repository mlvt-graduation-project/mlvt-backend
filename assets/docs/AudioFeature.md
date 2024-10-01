# API Documentation for Audio Features

## 1. Add a New Audio
- **API Endpoint**: `POST /audios`
- **Description**: Adds a new audio to the system. (Protected)
- **Input** (JSON body):
    ```json
    {
        "video_id": 123,
        "user_id": 456,
        "duration": 180,
        "lang": "en",
        "folder": "audios/2023/",
        "file_name": "audio.mp3"
    }
    ```
- **Response**:
    - `201 Created`: Audio added successfully.
    - `400 Bad Request`: Validation error.
    - `500 Internal Server Error`: Server-side issue.

## 2. Get Audio by ID
- **API Endpoint**: `GET /audios/{audio_id}`
- **Description**: Retrieves an audio by its ID and generates a presigned download URL. (Protected)
- **Input** (Path parameter):
    - `audio_id` (int): ID of the audio.
- **Response** (Example JSON response):
    ```json
    {
        "audio": {
            "id": 1,
            "video_id": 123,
            "user_id": 456,
            "duration": 180,
            "lang": "en",
            "folder": "audios/2023/",
            "file_name": "audio.mp3",
            "created_at": "2023-09-01T12:34:56Z",
            "updated_at": "2023-09-01T12:34:56Z"
        },
        "download_url": "https://s3.amazonaws.com/examplebucket/audios/2023/audio.mp3?presigned-url"
    }
    ```
    - `404 Not Found`: Audio not found.

## 3. Generate Presigned Upload URL for Audio
- **API Endpoint**: `POST /audios/generate-presigned-url`
- **Description**: Generates a presigned URL to upload an audio file. (Protected)
- **Input** (Query parameters):
    - `file_name` (string): Name of the audio file.
    - `file_type` (string): MIME type of the file (e.g., `audio/mpeg`).
- **Response** (Example JSON response):
    ```json
    {
        "upload_url": "https://s3.amazonaws.com/examplebucket/audios/2023/audio.mp3?presigned-url"
    }
    ```
    - `500 Internal Server Error`: Server-side issue.

## 4. Generate Presigned Download URL for Audio
- **API Endpoint**: `GET /audios/{audioID}/download-url`
- **Description**: Generates a presigned URL to download an audio file from S3. (Protected)
- **Input** (Path parameter):
    - `audioID` (int): ID of the audio.
- **Response** (Example JSON response):
    ```json
    {
        "download_url": "https://s3.amazonaws.com/examplebucket/audios/2023/audio.mp3?presigned-url"
    }
    ```
    - `500 Internal Server Error`: Server-side issue.

## 5. Get Audio by User and Audio ID
- **API Endpoint**: `GET /audios/{audioID}/user/{userID}`
- **Description**: Retrieves an audio for a specific user and generates a presigned download URL. (Protected)
- **Input** (Path parameters):
    - `audioID` (int): ID of the audio.
    - `userID` (int): ID of the user.
- **Response** (Example JSON response):
    ```json
    {
        "audio": {
            "id": 1,
            "video_id": 123,
            "user_id": 456,
            "duration": 180,
            "lang": "en",
            "folder": "audios/2023/",
            "file_name": "audio.mp3",
            "created_at": "2023-09-01T12:34:56Z",
            "updated_at": "2023-09-01T12:34:56Z"
        },
        "download_url": "https://s3.amazonaws.com/examplebucket/audios/2023/audio.mp3?presigned-url"
    }
    ```
    - `500 Internal Server Error`: Server-side issue.

## 6. Get Audio by Video and Audio ID
- **API Endpoint**: `GET /audios/{audioID}/video/{videoID}`
- **Description**: Retrieves an audio for a specific video and generates a presigned download URL. (Protected)
- **Input** (Path parameters):
    - `audioID` (int): ID of the audio.
    - `videoID` (int): ID of the video.
- **Response** (Example JSON response):
    ```json
    {
        "audio": {
            "id": 1,
            "video_id": 123,
            "user_id": 456,
            "duration": 180,
            "lang": "en",
            "folder": "audios/2023/",
            "file_name": "audio.mp3",
            "created_at": "2023-09-01T12:34:56Z",
            "updated_at": "2023-09-01T12:34:56Z"
        },
        "download_url": "https://s3.amazonaws.com/examplebucket/audios/2023/audio.mp3?presigned-url"
    }
    ```
    - `500 Internal Server Error`: Server-side issue.

## 7. List Audios by User ID
- **API Endpoint**: `GET /audios/user/{user_id}`
- **Description**: Lists all audio files uploaded by a specific user. (Protected)
- **Input** (Path parameter):
    - `user_id` (int): ID of the user.
- **Response** (Example JSON response):
    ```json
    {
        "audios": [
            {
                "id": 1,
                "video_id": 123,
                "user_id": 456,
                "duration": 180,
                "lang": "en",
                "folder": "audios/2023/",
                "file_name": "audio1.mp3",
                "created_at": "2023-09-01T12:34:56Z",
                "updated_at": "2023-09-01T12:34:56Z"
            },
            {
                "id": 2,
                "video_id": 124,
                "user_id": 456,
                "duration": 120,
                "lang": "es",
                "folder": "audios/2023/",
                "file_name": "audio2.mp3",
                "created_at": "2023-09-02T12:34:56Z",
                "updated_at": "2023-09-02T12:34:56Z"
            }
        ]
    }
    ```
    - `500 Internal Server Error`: Server-side issue.

## 8. List Audios by Video ID
- **API Endpoint**: `GET /audios/video/{video_id}`
- **Description**: Lists all audio files associated with a specific video. (Protected)
- **Input** (Path parameter):
    - `video_id` (int): ID of the video.
- **Response** (Example JSON response):
    ```json
    {
        "audios": [
            {
                "id": 1,
                "video_id": 123,
                "user_id": 456,
                "duration": 180,
                "lang": "en",
                "folder": "audios/2023/",
                "file_name": "audio1.mp3",
                "created_at": "2023-09-01T12:34:56Z",
                "updated_at": "2023-09-01T12:34:56Z"
            },
            {
                "id": 2,
                "video_id": 123,
                "user_id": 789,
                "duration": 240,
                "lang": "es",
                "folder": "audios/2023/",
                "file_name": "audio2.mp3",
                "created_at": "2023-09-02T12:34:56Z",
                "updated_at": "2023-09-02T12:34:56Z"
            }
        ]
    }
    ```
    - `500 Internal Server Error`: Server-side issue.

## 9. Delete Audio by ID
- **API Endpoint**: `DELETE /audios/{id}`
- **Description**: Deletes an audio by its ID. (Protected)
- **Input** (Path parameter):
    - `id` (int): ID of the audio.
- **Response**:
    - `200 OK`: Audio deleted successfully.
    - `400 Bad Request`: Invalid audio ID.
    - `500 Internal Server Error`: Server-side issue.

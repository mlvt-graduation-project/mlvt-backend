# API Documentation for Transcription Features

## 1. Add a New Transcription
- **API Endpoint**: `POST /transcriptions`
- **Description**: Adds a new transcription to the system. (Protected)
- **Input** (JSON body):
    ```json
    {
        "video_id": 100,
        "user_id": 200,
        "text": "This is the transcription text",
        "lang": "en",
        "folder": "transcriptions/2023/",
        "file_name": "transcription.json"
    }
    ```
- **Response**:
    - `201 Created`: Transcription added successfully.
    - `400 Bad Request`: Validation error.
    - `500 Internal Server Error`: Server-side issue.

## 2. Get Transcription by ID
- **API Endpoint**: `GET /transcriptions/{transcription_id}`
- **Description**: Retrieves a transcription by its ID and generates a presigned download URL. (Protected)
- **Input** (Path parameter):
    - `transcription_id` (int): ID of the transcription.
- **Response** (Example JSON response):
    ```json
    {
        "transcription": {
            "id": 1,
            "video_id": 100,
            "user_id": 200,
            "text": "This is the transcription text",
            "lang": "en",
            "folder": "transcriptions/2023/",
            "file_name": "transcription.json",
            "created_at": "2023-09-01T12:34:56Z",
            "updated_at": "2023-09-01T12:34:56Z"
        },
        "download_url": "https://s3.amazonaws.com/examplebucket/transcriptions/2023/transcription.json?presigned-url"
    }
    ```
    - `404 Not Found`: Transcription not found.

## 3. Get Transcription by User ID and Transcription ID
- **API Endpoint**: `GET /transcriptions/{transcriptionID}/user/{userID}`
- **Description**: Retrieves a transcription for a specific user and generates a presigned download URL. (Protected)
- **Input** (Path parameters):
    - `transcriptionID` (int): ID of the transcription.
    - `userID` (int): ID of the user.
- **Response** (Example JSON response):
    ```json
    {
        "transcription": {
            "id": 1,
            "video_id": 100,
            "user_id": 200,
            "text": "This is the transcription text",
            "lang": "en",
            "folder": "transcriptions/2023/",
            "file_name": "transcription.json",
            "created_at": "2023-09-01T12:34:56Z",
            "updated_at": "2023-09-01T12:34:56Z"
        },
        "download_url": "https://s3.amazonaws.com/examplebucket/transcriptions/2023/transcription.json?presigned-url"
    }
    ```
    - `500 Internal Server Error`: Server-side error.

## 4. Get Transcription by Video ID and Transcription ID
- **API Endpoint**: `GET /transcriptions/{transcriptionID}/video/{videoID}`
- **Description**: Retrieves a transcription for a specific video and generates a presigned download URL. (Protected)
- **Input** (Path parameters):
    - `transcriptionID` (int): ID of the transcription.
    - `videoID` (int): ID of the video.
- **Response** (Example JSON response):
    ```json
    {
        "transcription": {
            "id": 1,
            "video_id": 100,
            "user_id": 200,
            "text": "This is the transcription text",
            "lang": "en",
            "folder": "transcriptions/2023/",
            "file_name": "transcription.json",
            "created_at": "2023-09-01T12:34:56Z",
            "updated_at": "2023-09-01T12:34:56Z"
        },
        "download_url": "https://s3.amazonaws.com/examplebucket/transcriptions/2023/transcription.json?presigned-url"
    }
    ```
    - `500 Internal Server Error`: Server-side error.

## 5. List Transcriptions by User ID
- **API Endpoint**: `GET /transcriptions/user/{user_id}`
- **Description**: Lists all transcriptions belonging to a specific user. (Protected)
- **Input** (Path parameter):
    - `user_id` (int): ID of the user.
- **Response** (Example JSON response):
    ```json
    {
        "transcriptions": [
            {
                "id": 1,
                "video_id": 100,
                "user_id": 200,
                "text": "This is the transcription text",
                "lang": "en",
                "folder": "transcriptions/2023/",
                "file_name": "transcription1.json",
                "created_at": "2023-09-01T12:34:56Z",
                "updated_at": "2023-09-01T12:34:56Z"
            },
            {
                "id": 2,
                "video_id": 101,
                "user_id": 200,
                "text": "Another transcription text",
                "lang": "es",
                "folder": "transcriptions/2023/",
                "file_name": "transcription2.json",
                "created_at": "2023-09-02T12:34:56Z",
                "updated_at": "2023-09-02T12:34:56Z"
            }
        ]
    }
    ```
    - `500 Internal Server Error`: Server-side issue.

## 6. List Transcriptions by Video ID
- **API Endpoint**: `GET /transcriptions/video/{video_id}`
- **Description**: Lists all transcriptions belonging to a specific video. (Protected)
- **Input** (Path parameter):
    - `video_id` (int): ID of the video.
- **Response** (Example JSON response):
    ```json
    {
        "transcriptions": [
            {
                "id": 1,
                "video_id": 100,
                "user_id": 200,
                "text": "This is the transcription text",
                "lang": "en",
                "folder": "transcriptions/2023/",
                "file_name": "transcription1.json",
                "created_at": "2023-09-01T12:34:56Z",
                "updated_at": "2023-09-01T12:34:56Z"
            }
        ]
    }
    ```
    - `500 Internal Server Error`: Server-side issue.

## 7. Delete Transcription by ID
- **API Endpoint**: `DELETE /transcriptions/{transcription_id}`
- **Description**: Deletes a transcription by its ID. (Protected)
- **Input** (Path parameter):
    - `transcription_id` (int): ID of the transcription.
- **Response**:
    - `200 OK`: Transcription deleted successfully.
    - `500 Internal Server Error`: Server-side issue.

## 8. Generate Presigned Upload URL for Transcription
- **API Endpoint**: `POST /transcriptions/generate-upload-url`
- **Description**: Generates a presigned URL to upload a transcription file. (Protected)
- **Input** (Query parameters):
    - `file_name` (string): Name of the transcription file.
    - `file_type` (string): MIME type of the file (e.g., `application/json`).
- **Response** (Example JSON response):
    ```json
    {
        "upload_url": "https://s3.amazonaws.com/examplebucket/transcriptions/2023/transcription.json?presigned-url"
    }
    ```
    - `500 Internal Server Error`: Server-side issue.

## 9. Generate Presigned Download URL for Transcription
- **API Endpoint**: `GET /transcriptions/{transcription_id}/download-url`
- **Description**: Generates a presigned URL to download a transcription file. (Protected)
- **Input** (Path parameter):
    - `transcription_id` (int): ID of the transcription.
- **Response** (Example JSON response):
    ```json
    {
        "download_url": "https://s3.amazonaws.com/examplebucket/transcriptions/2023/transcription.json?presigned-url"
    }
    ```
    - `500 Internal Server Error`: Server-side issue.

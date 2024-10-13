# API Documentation for Video Features

## 1. Add a New Video
- **API Endpoint**: POST /videos/
- **Description**: Adds a new video to the system. (Protected)
- **Input** (JSON body):
  ```json
  {
      "title": "My Video Title",
      "duration": 300,
      "description": "A description of the video",
      "file_name": "video.mp4",
      "folder": "videos/2023/",
      "image": "thumbnail.jpg",
      "user_id": 123
  }
  ```
- **Response**:
  - 201 Created: Video added successfully.
  - 400 Bad Request: Validation error.
  - 500 Internal Server Error: Server-side issue.

## 2. Generate Presigned Upload URL for Video
- **API Endpoint**: POST /videos/generate-upload-url/video
- **Description**: Generates a presigned URL to upload a video file to S3. (Protected)
- **Input** (Query parameters):
  - `file_name` (string): The name of the video file.
  - `file_type` (string): The MIME type of the video file (e.g., video/mp4).
- **Response** (Example JSON response):
  ```json
  {
      "upload_url": "https://s3.amazonaws.com/examplebucket/videos/2023/video.mp4?presigned-url"
  }
  ```
  - 500 Internal Server Error: Server-side issue.

## 3. Generate Presigned Upload URL for Image
- **API Endpoint**: POST /videos/generate-upload-url/image
- **Description**: Generates a presigned URL to upload an image (e.g., video thumbnail) to S3. (Protected)
- **Input** (Query parameters):
  - `file_name` (string): The name of the image file.
  - `file_type` (string): The MIME type of the image file (e.g., image/jpeg).
- **Response** (Example JSON response):
  ```json
  {
      "upload_url": "https://s3.amazonaws.com/examplebucket/thumbnails/thumbnail.jpg?presigned-url"
  }
  ```
  - 500 Internal Server Error: Server-side issue.

## 4. Generate Presigned Download URL for Video
- **API Endpoint**: GET /videos/{video_id}/download-url/video
- **Description**: Generates a presigned URL to download a video file from S3. (Protected)
- **Input** (Path parameter):
  - `video_id` (int): ID of the video file.
- **Response** (Example JSON response):
  ```json
  {
      "video_download_url": "https://s3.amazonaws.com/examplebucket/videos/2023/video.mp4?presigned-url"
  }
  ```
  - 500 Internal Server Error: Server-side issue.

## 5. Generate Presigned Download URL for Image
- **API Endpoint**: GET /videos/{video_id}/download-url/image
- **Description**: Generates a presigned URL to download an image (e.g., thumbnail) from S3. (Protected)
- **Input** (Path parameter):
  - `video_id` (int): ID of the video file.
- **Response** (Example JSON response):
  ```json
  {
      "image_download_url": "https://s3.amazonaws.com/examplebucket/thumbnails/thumbnail.jpg?presigned-url"
  }
  ```
  - 500 Internal Server Error: Server-side issue.

## 6. Get Video by ID
- **API Endpoint**: GET /videos/{video_id}
- **Description**: Fetches a video by its ID and generates presigned URLs for the video and image. (Protected)
- **Input** (Path parameter):
  - `video_id` (int): ID of the video.
- **Response** (Example JSON response):
  ```json
  {
      "video": {
          "id": 1,
          "title": "My Video Title",
          "duration": 300,
          "description": "A description of the video",
          "file_name": "video.mp4",
          "folder": "videos/2023/",
          "image": "thumbnail.jpg",
          "user_id": 123,
          "created_at": "2023-09-01T12:34:56Z",
          "updated_at": "2023-09-01T12:34:56Z"
      },
      "video_url": "https://s3.amazonaws.com/examplebucket/videos/2023/video.mp4?presigned-url",
      "image_url": "https://s3.amazonaws.com/examplebucket/thumbnails/thumbnail.jpg?presigned-url"
  }
  ```
  - 404 Not Found: Video not found.

## 7. Delete Video by ID
- **API Endpoint**: DELETE /videos/{video_id}
- **Description**: Deletes a video by its ID from the system. (Protected)
- **Input** (Path parameter):
  - `video_id` (int): ID of the video.
- **Response**:
  - 200 OK: Video deleted successfully.
  - 500 Internal Server Error: Server-side issue.

## 8. List Videos by User ID
- **API Endpoint**: GET /videos/user/{user_id}
- **Description**: Fetches all videos uploaded by a specific user along with presigned image URLs. (Protected)
- **Input** (Path parameter):
  - `user_id` (int): ID of the user.
- **Response** (Example JSON response):
  ```json
  {
      "videos": [
          {
              "id": 1,
              "title": "My First Video",
              "duration": 300,
              "description": "A description of the video",
              "file_name": "video1.mp4",
              "folder": "videos/2023/",
              "image": "thumbnail1.jpg",
              "user_id": 123,
              "created_at": "2023-09-01T12:34:56Z",
              "updated_at": "2023-09-01T12:34:56Z"
          },
          {
              "id": 2,
              "title": "My Second Video",
              "duration": 240,
              "description": "Another video description",
              "file_name": "video2.mp4",
              "folder": "videos/2023/",
              "image": "thumbnail2.jpg",
              "user_id": 123,
              "created_at": "2023-09-02T12:34:56Z",
              "updated_at": "2023-09-02T12:34:56Z"
          }
      ]
  }
  ```
  - 500 Internal Server Error: Server-side error.

## 9. Get Video Status
- **API Endpoint**: GET /videos/{video_id}/status
- **Description**: Retrieves the status of a specific video by its ID. (Protected)
- **Input** (Path parameter):
  - `video_id` (uint64): Video ID.
- **Response**:
  - 200 OK: Status of the video (e.g., `{"status": "success"}`).
  - 400 Bad Request: Invalid video ID.
  - 404 Not Found: Video not found.
  - 500 Internal Server Error: Server-side issue.

## 10. Update Video Status
- **API Endpoint**: PUT /videos/{video_id}/status
- **Description**: Updates the status of a specific video by its ID. (Protected)
- **Input**:
  - **Path parameter**:
    - `video_id` (uint64): Video ID.
  - **Body (JSON)**:
    ```json
    {
        "status": "processing"
    }
    ```
    - Status must be one of: `raw`, `processing`, `failed`, `success`.
- **Response**:
  - 200 OK: Status updated successfully.
  - 400 Bad Request: Invalid input.
  - 404 Not Found: Video not found.
  - 500 Internal Server Error: Server-side issue.
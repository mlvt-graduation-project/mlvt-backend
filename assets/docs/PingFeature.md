# Ping Status API Guide

## Endpoints

### 1. Speech-to-Text

- **URL:** `/speech-to-text/{id}`
- **Method:** `GET`
- **Description:** Retrieves the status of a speech-to-text transcription task.
- **Parameters:**
  - `id` (path) – The unique identifier of the transcription task.

- **Response:**
  - **200 OK**
    ```json
    {
      "status": "completed"
    }
    ```
  - **400 Bad Request**
    ```json
    {
      "error": "Invalid ID"
    }
    ```
  - **500 Internal Server Error**
    ```json
    {
      "error": "Error message"
    }
    ```


### 2. Text-to-Text

- **URL:** `/text-to-text/{id}`
- **Method:** `GET`
- **Description:** Retrieves the status of a text-to-text transcription task.
- **Parameters:**
- `id` (path) – The unique identifier of the transcription task.

- **Response:** Same as **Speech-to-Text**.


### 3. Text-to-Speech

- **URL:** `/text-to-speech/{id}`
- **Method:** `GET`
- **Description:** Retrieves the status of a text-to-speech audio task.
- **Parameters:**
- `id` (path) – The unique identifier of the audio task.

- **Response:** Same as **Speech-to-Text**.


### 4. Voice Cloning

- **URL:** `/voice-cloning/{id}`
- **Method:** `GET`
- **Description:** Retrieves the status of a voice cloning audio task.
- **Parameters:**
- `id` (path) – The unique identifier of the audio task.

- **Response:** Same as **Speech-to-Text**.


### 5. Lip Sync

- **URL:** `/lipsync/{id}`
- **Method:** `GET`
- **Description:** Retrieves the status of a lip-sync video task.
- **Parameters:**
- `id` (path) – The unique identifier of the video task.

- **Response:** Same as **Speech-to-Text**.


### 6. Full Pipeline

- **URL:** `/full-pipeline/{id}`
- **Method:** `GET`
- **Description:** Retrieves the status of a full pipeline video task.
- **Parameters:**
- `id` (path) – The unique identifier of the video task.

- **Response:** Same as **Speech-to-Text**.


## Response Structure

All successful responses will have the following JSON structure:

```json
{
"status": "raw" | "processing" | "succeeded" | "failed"
}

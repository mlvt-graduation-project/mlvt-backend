## Front-End Developer Guidelines

This guide provides essential information for Front-End (FE) developers to understand how user avatars and videos (along with their frames) are managed within the project. It includes details about file locations, naming conventions, default behaviors, and data relationships.

### Table of Contents

- [Front-End Developer Guidelines](#front-end-developer-guidelines)
  - [Table of Contents](#table-of-contents)
  - [User Avatars](#user-avatars)
  - [Videos and Video Frames](#videos-and-video-frames)
  - [Default Password](#default-password)
  - [Uploading to AWS S3](#uploading-to-aws-s3)
  - [User ID Assignment](#user-id-assignment)
  - [Additional Information](#additional-information)

---

### User Avatars

- **File Location:** `assets/avatars`
  
- **Processing:**
  - The seeder loops through all avatar images in the `assets/avatars` directory.
  - For each avatar image:
    - Creates a new user in the database.
    - Uploads the avatar image to AWS S3.
    - Sets the default password to `"capigiba"` for the user.

- **Naming Convention:**
  - Avatar files retain their original names.
  - **Example:** `john_doe.jpg` should be located at `assets/avatars/john_doe.jpg`.

### Videos and Video Frames

- **File Location:** `assets/videos`
  
- **Processing:**
  - The seeder loops through all video files in the `assets/videos` directory.
  - For each video file:
    - Identifies the corresponding frame image based on naming convention (e.g., `capi.mp4` has `capi_thumbnail.jpg`).
    - Uploads both the video file and its frame image to AWS S3.
    - Creates a video record in the database.
    - Assigns the video to a user, selected randomly from the existing users in the database.

- **Naming Convention:**
  - **Video File:** `<video_filename>.<extension>` (e.g., `capi.mp4`)
  - **Frame Image:** `<video_filename>_thumbnail.jpg` (e.g., `capi_thumbnail.jpg`)
    - **Supported Extensions for Frames:** Currently only `.jpg` is supported.

- **Example:**
  - **Video File:** `capi.mp4`
  - **Frame Image:** `capi_thumbnail.jpg`

### Default Password

- **Password:** All users created via the seeder have the default password `"capigiba"`.
  
- **Usage:**
  - This password is intended for development and testing purposes.
  - Ensure that this password is changed in production environments to maintain security.

### Uploading to AWS S3

- **Avatars:**
  - **Destination Folder on S3:** Defined by the environment variable `AWSAvatarFolder`.
  - **Process:**
    1. The seeder reads each avatar image from `assets/avatars`.
    2. Each image is uploaded to the specified S3 folder.
    3. The user's record is updated with the S3 file name and folder details.

- **Videos and Frames:**
  - **Destination Folders on S3:**
    - **Videos:** Defined by the environment variable `AWSVideosFolder`.
    - **Frames (Thumbnails):** Defined by the environment variable `AWSVideoFramesFolder`.
  - **Process:**
    1. The seeder reads each video file from `assets/videos`.
    2. The corresponding frame image is identified based on the naming convention.
    3. Both the video file and its frame image are uploaded to their respective S3 folders.
    4. The video's record is updated with the S3 file names and folder details.
    5. The video is assigned to a user randomly selected from the existing users in the database.

**Note:** Ensure that AWS credentials and permissions are correctly configured to allow upload and deletion operations on the specified S3 buckets.

### User ID Assignment

- **Random Assignment:**
  - When seeding videos, each video is randomly assigned to an existing user in the database.
  
- **Process:**
  1. The seeder selects a user ID randomly from the list of available user IDs.
  2. The video record is associated with the selected user ID in the database.

- **Implications:**
  - This random assignment ensures an even distribution of videos across users.
  - FE developers can retrieve videos based on user IDs for display purposes.

### Additional Information

- **Supported File Types:**
  - **Avatars:** `.jpg`, `.jpeg`, `.png`, `.gif`
  - **Videos:** `.mp4`, `.avi`, `.mov`, `.mkv`
  - **Frames (Thumbnails):** `.jpg` (currently supported)

- **Seeding Process:**
  1. **Avatars:**
     - Loop through all avatar images in `assets/avatars`.
     - For each avatar:
       - Create a user with default password `"capigiba"`.
       - Upload the avatar to S3.
       - Update the user's avatar information in the database.
  
  2. **Videos:**
     - Loop through all video files in `assets/videos`.
     - For each video:
       - Identify the corresponding frame image.
       - Upload both video and frame to S3.
       - Create a video record in the database.
       - Assign the video to a random user.

- **Data Cleanup:**
  - A separate cleaner service handles the deletion of seeded data from both the database and AWS S3.
  - It performs the following:
    - Deletes avatar images from S3.
    - Deletes video files and frame images from S3.
    - Removes user records and video records from the database.

- **Password Security:**
  - While the default password `"capigiba"` is suitable for development and testing, ensure that user passwords are handled securely in production environments.
  - Encourage users to change their passwords upon first login.

- **Error Handling:**
  - The seeder logs errors encountered during the seeding process but continues with subsequent operations.
  - FE developers should be aware of potential inconsistencies and implement appropriate error handling mechanisms.
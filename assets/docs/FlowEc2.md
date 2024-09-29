```mermaid
sequenceDiagram
    participant Client
    participant Server as Server 
    participant API_Gateway as API Gateway
    participant EC2_Instance as EC2 Instance
    participant Model_Processing as Model Processing

    Client->>Server: Request xử lý video
    Server->>API_Gateway: Gửi API request (HTTP/HTTPS)
    API_Gateway->>EC2_Instance: Chuyển tiếp request
    EC2_Instance->>Model_Processing: Khởi động mô hình xử lý (CUDA/GPU)
    Model_Processing->>EC2_Instance: Trả kết quả xử lý
    EC2_Instance->>API_Gateway: Gửi kết quả xử lý về
    API_Gateway->>Server: Chuyển kết quả về Server
    Server->>Client: Trả kết quả cuối cùng cho Client

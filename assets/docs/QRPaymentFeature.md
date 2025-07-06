# QR Payment Feature

This document describes the QR payment feature implementation using the VietQR.io API for the MLVT backend.

## Overview

The QR payment feature allows users to purchase tokens by generating VietQR codes that can be scanned with any Vietnamese banking app. The system integrates with VietQR.io API and uses a VietinBank account to receive payments.

## Features

- **Multiple Payment Options**: 5 predefined packages from 500k VND to 10M VND
- **Automatic QR Generation**: Creates QR codes with unique transaction IDs
- **Payment Tracking**: Monitors payment status and transaction history
- **Automatic Token Credit**: Adds tokens to user wallet upon payment confirmation
- **Admin Monitoring**: Admin panel to view pending payments
- **Transaction History**: Users can view their payment history

## Architecture

### Entities

- **PaymentTransaction**: Main payment entity stored in MongoDB
- **PaymentOption**: Predefined payment packages (500k, 1m, 2m, 5m, 10m)
- **VietQRRequest/Response**: API communication structures

### Components

1. **payment_entity.go**: Data models and payment options
2. **payment_repo.go**: MongoDB repository for payment operations
3. **payment_service.go**: Business logic and VietQR API integration
4. **payment_handler.go**: HTTP handlers for REST API endpoints

## API Endpoints

### Public Endpoints

- `GET /api/payment/options` - Get available payment options

### User Endpoints (Authenticated)

- `POST /api/payment/create?user_id=123&option=1m` - Create payment QR
- `GET /api/payment/{payment_id}` - Get payment by ID
- `GET /api/payment/transaction/{transaction_id}` - Get payment by transaction ID
- `GET /api/payment/user-payments?user_id=123` - Get user's payment history
- `POST /api/payment/confirm/{transaction_id}` - Confirm payment (adds tokens)
- `POST /api/payment/cancel/{payment_id}` - Cancel pending payment

### Admin Endpoints

- `GET /api/payment/pending` - Get all pending payments

## Payment Flow

1. **User initiates payment**:
   - User selects payment option (500k, 1m, 2m, 5m, 10m)
   - System generates unique transaction ID
   - Calls VietQR API to create QR code
   - Returns QR code and payment details

2. **User scans QR and pays**:
   - User scans QR with banking app
   - Makes payment to VietinBank account
   - Bank transfer includes transaction ID in message

3. **Payment confirmation**:
   - Admin or automated system confirms payment
   - System adds tokens to user wallet
   - Payment status updated to "completed"

## Configuration

### Environment Variables

Add these to your `.env` file:

```env
# VietQR API credentials (get from https://my.vietqr.io)
VIETQR_CLIENT_ID=your_client_id_here
VIETQR_API_KEY=your_api_key_here

# VietinBank account information
VIETINBANK_ACCOUNT_NO=your_account_number
VIETINBANK_ACCOUNT_NAME=YOUR ACCOUNT NAME
VIETINBANK_BIN_CODE=970415
```

### Payment Options

| Option | Tokens | VND Price | Rate |
|--------|--------|-----------|------|
| 500k   | 500    | 500,000   | 1:1000 |
| 1m     | 1,000  | 1,000,000 | 1:1000 |
| 2m     | 2,000  | 2,000,000 | 1:1000 |
| 5m     | 5,000  | 5,000,000 | 1:1000 |
| 10m    | 10,000 | 10,000,000| 1:1000 |

## Example Usage

### Create Payment QR

```bash
curl -X POST "http://localhost:8080/api/payment/create?user_id=1&option=1m" \
  -H "Authorization: Bearer your_jwt_token"
```

Response:
```json
{
  "id": "64a7b8c9d1e2f3a4b5c6d7e8",
  "user_id": 1,
  "transaction_id": "TXN_1703123456789012345",
  "payment_option": "1m",
  "token_amount": 1000,
  "vnd_amount": 1000000,
  "status": "pending",
  "qr_code": "00020101021238560010A0000007270126...",
  "qr_data_url": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...",
  "created_at": "2023-12-21T10:30:56Z"
}
```

### Check Payment Status

```bash
curl -X GET "http://localhost:8080/api/payment/transaction/TXN_1703123456789012345" \
  -H "Authorization: Bearer your_jwt_token"
```

### Confirm Payment (Admin)

```bash
curl -X POST "http://localhost:8080/api/payment/confirm/TXN_1703123456789012345" \
  -H "Authorization: Bearer admin_jwt_token"
```

## Database Schema

### PaymentTransaction (MongoDB)

```json
{
  "_id": "ObjectId",
  "user_id": "uint64",
  "transaction_id": "string",
  "payment_option": "string",
  "token_amount": "int64",
  "vnd_amount": "int64",
  "status": "string",
  "qr_code": "string",
  "qr_data_url": "string",
  "created_at": "datetime",
  "updated_at": "datetime",
  "completed_at": "datetime"
}
```

## Security Considerations

1. **API Keys**: Store VietQR credentials securely in environment variables
2. **Authentication**: All endpoints require user authentication
3. **Transaction IDs**: Use unique, non-guessable transaction identifiers
4. **Payment Verification**: Manual confirmation prevents fraudulent token credits
5. **Rate Limiting**: Consider adding rate limits for payment creation

## Monitoring and Analytics

- Track payment success rates
- Monitor pending payment queue
- Analyze popular payment options
- Generate revenue reports

## Future Enhancements

1. **Webhook Integration**: Automatic payment confirmation via bank webhooks
2. **Multiple Banks**: Support for other Vietnamese banks
3. **Dynamic Pricing**: Flexible token rates and promotional pricing
4. **Payment Expiry**: Auto-cancel payments after timeout
5. **Refund System**: Handle payment refunds and cancellations

## Troubleshooting

### Common Issues

1. **QR Generation Fails**:
   - Check VietQR API credentials
   - Verify account information format
   - Check API rate limits

2. **Payment Not Confirmed**:
   - Verify transaction ID in bank transfer
   - Check payment status in admin panel
   - Ensure sufficient account permissions

3. **Tokens Not Added**:
   - Check wallet service integration
   - Verify user ID mapping
   - Review transaction logs

### Error Codes

- `invalid_payment_option`: Unknown payment package
- `payment_not_found`: Transaction ID not found
- `payment_not_pending`: Cannot confirm non-pending payment
- `failed_to_generate_qr`: VietQR API error
- `insufficient_balance`: User wallet error 
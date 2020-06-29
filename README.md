# setel-orderPayment-backend
Simple order and payment backend service implementation

## Lambdas and endpoints
Sample running lambda services are hosted on AWS

Base service URL: `https://fymcvrpts6.execute-api.ap-southeast-1.amazonaws.com/`

### Create Order API
Create Order API - POST `/order`

Sample request body:
```
{
    "Description": "Hello World",
    "Amount": 99.99
}
```

### Get Order API
Get Order(s) API - GET `/order` or `/order/{id}`

### Update Order API
Update Order API - PATCH `/order/{id}`

Sample request body:
```
{
    "Description": "Hello World",
    "Amount": 99.99
    "status": 1
}
```

### Create Payment API
Create Payment API - POST `/payment`

Sample request body:
```
{
    "AuthDetails": "auth data...",
    "OrderDetails": {
        "id": "abc123",
        "status": "Created",
        "description": "Test",
        "amount": 9.99
    },
    "ManualTrigger": true
}
```

### Get Payment API
Get Payments API - GET `/payment`
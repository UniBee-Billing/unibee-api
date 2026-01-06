# Embedded Frontend (Multi-Gateway Payment Integration)

This directory contains the frontend payment code for multiple payment gateways using their respective Embedded APIs (Stripe, PayPal, etc.).

## Development

```bash
# Start development server
make frontend.dev
# or
cd embedded && npm start
```

## Build

```bash
# Build frontend code and copy to resource/embedded
make frontend.build

# Clean build files
make frontend.clean

# Build everything (frontend + Go backend)
make build.all
```

## Build Process

1. Frontend code is developed in the `embedded/` directory
2. When running `make frontend.build`:
   - Automatically installs dependencies
   - Builds React application
   - Copies build artifacts to `resource/embedded/` directory
3. Go server can directly access static files in `resource/embedded/`

## File Structure

```
embedded/           # Frontend source code directory
├── src/           # React source code
│   ├── App.jsx    # Main React component with routing
│   ├── components/ # Payment gateway components
│   │   ├── StripeCheckout.jsx       # Stripe Embedded Checkout
│   │   ├── PayPalCheckout.jsx       # PayPal Embedded Checkout
│   │   └── BlockonomicsCheckout.jsx # Blockonomics Crypto Checkout
│   ├── services/  # Payment service for API integration
│   │   └── paymentService.js     # Unified payment service
│   └── test.html  # Test page for development
├── public/        # Static assets
└── package.json   # Dependencies configuration

resource/embedded/  # Build output directory (accessed by Go server)
├── index.html     # Main page
├── static/        # Static assets
└── asset-manifest.json
```

## Integration with Go Backend

The frontend is integrated with the Go backend system:

1. **Payment Detail API**: Calls `/system/payment/detail?paymentId=xxxx` to get payment information
2. **Action Data**: Extracts gateway-specific configuration from the `action` field in the payment response
3. **Multi-Gateway Support**: Supports both Stripe and PayPal embedded checkouts
4. **Dynamic Configuration**: Uses gateway-specific keys and secrets from backend configuration

### Stripe Integration
- **Client Secret**: Uses `stripeClientSecret` from action data for Stripe Embedded Checkout
- **Publishable Key**: Uses `stripeAPIKey` from action data for Stripe initialization

### PayPal Integration
- **Order ID**: Uses `paypalOrderID` from action data for PayPal order creation
- **Client ID**: Uses `paypalClientId` from action data for PayPal SDK initialization

### Blockonomics Integration
- **Payment Address**: Uses `blockonomicsAddress` from action data for crypto payment address
- **API Key**: Uses `blockonomicsApiKey` from action data for API authentication
- **Order ID**: Uses `blockonomicsOrderId` from action data for order tracking
- **Amount**: Uses `blockonomicsAmount` from action data for payment amount
- **Currency**: Uses `blockonomicsCurrency` from action data for crypto currency type

## Usage

### Development
```bash
# Start development server
make frontend.dev
# or
cd embedded && npm start
```

### Production
```bash
# Build and deploy
make frontend.build
```

### Testing
1. Start your Go backend server
2. Create a payment record in your system
3. Access specific payment gateway:
   - Stripe: `http://localhost:8088/embedded/stripe?paymentId=YOUR_PAYMENT_ID`
   - PayPal: `http://localhost:8088/embedded/paypal?paymentId=YOUR_PAYMENT_ID`
   - Blockonomics: `http://localhost:8088/embedded/blockonomics?paymentId=YOUR_PAYMENT_ID`
   - Alipay: `http://localhost:8088/embedded/alipay?paymentId=YOUR_PAYMENT_ID` (future)
4. Or use the test page: `http://localhost:8088/embedded/test.html`

## API Integration

The frontend expects the following API response structure from `/system/payment/detail`:

```json
{
  "paymentDetail": {
    "payment": {
      "action": {
        "stripeSessionId": "cs_xxx",
        "stripeClientSecret": "cs_xxx_secret_xxx",
        "stripeAPIKey": "pk_test_xxx",
        "paypalOrderID": "PAYPAL_ORDER_ID",
        "paypalClientId": "PAYPAL_CLIENT_ID",
        "blockonomicsAddress": "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
        "blockonomicsApiKey": "BLOCKONOMICS_API_KEY",
        "blockonomicsOrderId": "BLOCKONOMICS_ORDER_ID",
        "blockonomicsAmount": "0.001",
        "blockonomicsCurrency": "BTC"
      }
    },
    "gateway": {
      "gatewayKey": "pk_test_xxx",
      "gatewaySecret": "sk_test_xxx"
    }
  }
}
```

### Action Data Structure

The `action` field contains gateway-specific payment configuration:

#### Stripe Configuration
- **stripeSessionId**: Stripe Checkout Session ID
- **stripeClientSecret**: Stripe Client Secret for Embedded Checkout
- **stripeAPIKey**: Stripe Publishable Key

#### PayPal Configuration
- **paypalOrderID**: PayPal Order ID for order creation
- **paypalClientId**: PayPal Client ID for SDK initialization

#### Blockonomics Configuration
- **blockonomicsAddress**: Cryptocurrency payment address
- **blockonomicsApiKey**: Blockonomics API key for authentication
- **blockonomicsOrderId**: Blockonomics order ID for tracking
- **blockonomicsAmount**: Payment amount in cryptocurrency
- **blockonomicsCurrency**: Cryptocurrency type (BTC, USDT, etc.)

### Payment Status

Payment status is obtained directly from the `/system/payment/detail` API response, eliminating the need for separate session status checks.

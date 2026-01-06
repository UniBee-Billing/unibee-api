# UniBee Payment iframe Integration Guide

## Overview

UniBee Payment supports iframe integration for seamless payment experiences. By setting `PaymentUIMode` to `"embedded"` or `"custom"` in API requests, you can embed payment pages directly into your application using iframe, eliminating the need for page redirects.

## API Integration

### Supported Endpoints

The following API endpoints support iframe integration:

| Endpoint | Description | PaymentUIMode Support |
|----------|-------------|----------------------|
| `/checkout/subscription/create_submit` | Create subscription (checkout) | ✅ |
| `/merchant/subscription/create_submit` | Create subscription (merchant) | ✅ |
| `/merchant/subscription/update_submit` | Update subscription (merchant) | ✅ |
| `/merchant/payment/new` | Create one-time payment | ✅ |

### Request Parameters

Add the `PaymentUIMode` parameter to your API requests:

```json
{
  "paymentUIMode": "embedded",  // or "custom"
  // ... other parameters
}
```

### Response Structure

When `PaymentUIMode` is set to `"embedded"` or `"custom"`, the API returns:

```json
{
  "status": 10,
  "paymentId": "pay_xxx",
  "invoiceId": "inv_xxx", 
  "link": "https://api.unibee.dev/embedded/payment_checker?paymentId=pay_xxx&env=prod",
  "action": {
    // Payment gateway specific data
  }
}
```

**Key Fields**:
- `link`: The payment page URL to embed in iframe
- `paymentId`: Unique payment identifier
- `invoiceId`: Unique invoice identifier
- `action`: Gateway-specific payment data

## Frontend Integration

### 1. Basic iframe Setup

```html
<!DOCTYPE html>
<html>
<head>
  <title>Payment Integration</title>
  <style>
    .payment-modal {
      display: none;
      position: fixed;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background: rgba(0,0,0,0.5);
      z-index: 1000;
    }
    
    .payment-content {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      background: white;
      border-radius: 8px;
      width: 90%;
      max-width: 800px;
      height: 80%;
      max-height: 600px;
    }
    
    .payment-iframe {
      width: 100%;
      height: 100%;
      border: none;
    }
  </style>
</head>
<body>
  <button onclick="openPayment()">Open Payment</button>
  
  <div id="paymentModal" class="payment-modal">
    <div class="payment-content">
      <iframe id="paymentIframe" class="payment-iframe" src="about:blank"></iframe>
    </div>
  </div>

  <script>
    // Message listener for iframe communication
    window.addEventListener('message', function(event) {
      console.log('Received message from iframe:', event.data);
      
      // Security verification (specify exact domain in production)
      // if (event.origin !== 'https://api.unibee.dev') return;
      
      const { type, paymentId, invoiceId } = event.data;
      
      switch (type) {
        case 'UniBee_PaymentSuccess':
          handlePaymentSuccess(paymentId, invoiceId);
          break;
        case 'UniBee_PaymentFailed':
          handlePaymentFailed(paymentId, invoiceId);
          break;
        case 'UniBee_PaymentCancelled':
          handlePaymentCancelled(paymentId, invoiceId);
          break;
      }
    });
    
    function openPayment() {
      // Call your API to create payment
      createPayment()
        .then(response => {
          const paymentUrl = response.link;
          document.getElementById('paymentIframe').src = paymentUrl;
          document.getElementById('paymentModal').style.display = 'block';
        })
        .catch(error => {
          console.error('Failed to create payment:', error);
        });
    }
    
    function closePayment() {
      document.getElementById('paymentModal').style.display = 'none';
      document.getElementById('paymentIframe').src = 'about:blank';
    }
    
    function handlePaymentSuccess(paymentId, invoiceId) {
      console.log('Payment successful:', { paymentId, invoiceId });
      closePayment();
      // Add your success handling logic here
      alert('Payment completed successfully!');
    }
    
    function handlePaymentFailed(paymentId, invoiceId) {
      console.log('Payment failed:', { paymentId, invoiceId });
      closePayment();
      // Add your failure handling logic here
      alert('Payment failed. Please try again.');
    }
    
    function handlePaymentCancelled(paymentId, invoiceId) {
      console.log('Payment cancelled:', { paymentId, invoiceId });
      closePayment();
      // Add your cancellation handling logic here
      alert('Payment was cancelled.');
    }
    
    // API call example
    async function createPayment() {
      const response = await fetch('/api/merchant/payment/new', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer your-token'
        },
        body: JSON.stringify({
          paymentUIMode: 'embedded',
          gatewayId: 1,
          amount: 1000,
          currency: 'USD',
          // ... other required parameters
        })
      });
      
      return await response.json();
    }
  </script>
</body>
</html>
```

### 2. React Integration Example

```jsx
import React, { useState, useEffect } from 'react';

const PaymentModal = ({ paymentUrl, onClose, onSuccess, onFailure, onCancel }) => {
  useEffect(() => {
    const handleMessage = (event) => {
      // Security verification
      // if (event.origin !== 'https://api.unibee.dev') return;
      
      const { type, paymentId, invoiceId } = event.data;
      
      switch (type) {
        case 'UniBee_PaymentSuccess':
          onSuccess?.(paymentId, invoiceId);
          break;
        case 'UniBee_PaymentFailed':
          onFailure?.(paymentId, invoiceId);
          break;
        case 'UniBee_PaymentCancelled':
          onCancel?.(paymentId, invoiceId);
          break;
      }
    };
    
    window.addEventListener('message', handleMessage);
    return () => window.removeEventListener('message', handleMessage);
  }, [onSuccess, onFailure, onCancel]);
  
  return (
    <div className="payment-modal">
      <div className="payment-content">
        <iframe 
          src={paymentUrl} 
          className="payment-iframe"
          title="Payment"
        />
      </div>
    </div>
  );
};

const PaymentPage = () => {
  const [showPayment, setShowPayment] = useState(false);
  const [paymentUrl, setPaymentUrl] = useState('');
  
  const createPayment = async () => {
    try {
      const response = await fetch('/api/merchant/payment/new', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer your-token'
        },
        body: JSON.stringify({
          paymentUIMode: 'embedded',
          gatewayId: 1,
          amount: 1000,
          currency: 'USD',
        })
      });
      
      const data = await response.json();
      setPaymentUrl(data.link);
      setShowPayment(true);
    } catch (error) {
      console.error('Failed to create payment:', error);
    }
  };
  
  const handleSuccess = (paymentId, invoiceId) => {
    setShowPayment(false);
    console.log('Payment successful:', { paymentId, invoiceId });
    // Handle success
  };
  
  const handleFailure = (paymentId, invoiceId) => {
    setShowPayment(false);
    console.log('Payment failed:', { paymentId, invoiceId });
    // Handle failure
  };
  
  const handleCancel = (paymentId, invoiceId) => {
    setShowPayment(false);
    console.log('Payment cancelled:', { paymentId, invoiceId });
    // Handle cancellation
  };
  
  return (
    <div>
      <button onClick={createPayment}>Create Payment</button>
      
      {showPayment && (
        <PaymentModal
          paymentUrl={paymentUrl}
          onClose={() => setShowPayment(false)}
          onSuccess={handleSuccess}
          onFailure={handleFailure}
          onCancel={handleCancel}
        />
      )}
    </div>
  );
};
```

## Message Communication Protocol

### Message Format

The iframe sends messages to the parent window in the following format:

```javascript
{
  type: 'UniBee_PaymentSuccess',    // or PaymentFailed, PaymentCancelled
  paymentId: 'pay_xxx',             // Payment ID
  invoiceId: 'inv_xxx'              // Invoice ID (optional)
}
```

### Message Types

| Type | Description | Trigger |
|------|-------------|---------|
| `UniBee_PaymentSuccess` | Payment successful | When payment status changes to success |
| `UniBee_PaymentFailed` | Payment failed | When payment status changes to failed |
| `UniBee_PaymentCancelled` | Payment cancelled | When payment status changes to cancelled |

## Complete Integration Flow

### 1. Create Payment

```javascript
// Step 1: Call API to create payment
const response = await fetch('/api/merchant/payment/new', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer your-token'
  },
  body: JSON.stringify({
    paymentUIMode: 'embedded',  // Enable iframe mode
    gatewayId: 1,
    amount: 1000,
    currency: 'USD',
    // ... other parameters
  })
});

const paymentData = await response.json();
```

### 2. Embed Payment Page

```javascript
// Step 2: Get payment URL from response
const paymentUrl = paymentData.link; // e.g., "https://api.unibee.dev/embedded/payment_checker?paymentId=pay_xxx&env=prod"

// Step 3: Create iframe and load payment page
const iframe = document.createElement('iframe');
iframe.src = paymentUrl;
iframe.style.width = '100%';
iframe.style.height = '600px';
iframe.style.border = 'none';

document.getElementById('payment-container').appendChild(iframe);
```

### 3. Listen for Messages

```javascript
// Step 4: Listen for payment status messages
window.addEventListener('message', function(event) {
  // Security verification
  if (event.origin !== 'https://api.unibee.dev') return;
  
  const { type, paymentId, invoiceId } = event.data;
  
  switch (type) {
    case 'UniBee_PaymentSuccess':
      // Payment completed successfully
      console.log('Payment successful:', paymentId);
      // Close iframe, redirect, or update UI
      break;
      
    case 'UniBee_PaymentFailed':
      // Payment failed
      console.log('Payment failed:', paymentId);
      // Show error message, retry option, etc.
      break;
      
    case 'UniBee_PaymentCancelled':
      // Payment was cancelled by user
      console.log('Payment cancelled:', paymentId);
      // Close iframe, show cancellation message
      break;
  }
});
```

## Security Considerations

### 1. Message Origin Verification

```javascript
window.addEventListener('message', function(event) {
  // Always verify message origin in production
  if (event.origin !== 'https://api.unibee.dev') {
    console.warn('Ignoring message from untrusted origin:', event.origin);
    return;
  }
  
  // Process message...
});
```

### 2. Message Content Validation

```javascript
function isValidMessage(data) {
  return data && 
         typeof data.type === 'string' && 
         data.type.startsWith('UniBee_') &&
         typeof data.paymentId === 'string';
}

window.addEventListener('message', function(event) {
  if (!isValidMessage(event.data)) {
    console.warn('Invalid message format:', event.data);
    return;
  }
  
  // Process message...
});
```

## Best Practices

### 1. Responsive Design

```css
.payment-iframe {
  width: 100%;
  height: 600px;
  border: none;
  border-radius: 8px;
}

@media (max-width: 768px) {
  .payment-iframe {
    height: 500px;
  }
}
```

### 2. Loading States

```javascript
function showPaymentModal(paymentUrl) {
  const modal = document.getElementById('paymentModal');
  const iframe = document.getElementById('paymentIframe');
  
  // Show loading state
  modal.style.display = 'block';
  iframe.style.display = 'none';
  
  // Load payment page
  iframe.onload = () => {
    iframe.style.display = 'block';
    // Hide loading state
  };
  
  iframe.src = paymentUrl;
}
```

### 3. Timeout Handling

```javascript
let paymentTimeout;

function openPayment() {
  // Set timeout (e.g., 30 minutes)
  paymentTimeout = setTimeout(() => {
    closePayment();
    alert('Payment timeout. Please try again.');
  }, 30 * 60 * 1000);
  
  // Open payment...
}

function closePayment() {
  if (paymentTimeout) {
    clearTimeout(paymentTimeout);
    paymentTimeout = null;
  }
  // Close modal...
}
```

## API Reference

### PaymentUIMode Values

| Value | Description |
|-------|-------------|
| `"hosted"` | Redirect to hosted payment page (default) |
| `"embedded"` | Return URL for iframe embedding |
| `"custom"` | Return URL for custom integration |

### Required Headers

```javascript
{
  'Content-Type': 'application/json',
  'Authorization': 'Bearer YOUR-API-KEY'
}
```

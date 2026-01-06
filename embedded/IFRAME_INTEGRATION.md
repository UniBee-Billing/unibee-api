# UniBee Embedded Payment iframe Integration Guide

## Overview

UniBee Embedded Payment supports running in iframe and communicates with parent window through `postMessage` API. When payment status changes, the iframe sends messages to the parent window, which can close the modal and handle payment results accordingly.

## Supported Payment Methods

- **Stripe**: Credit card payment
- **PayPal**: PayPal payment
- **Blockonomics**: Cryptocurrency payment (BTC/USDT)

## Page Paths

All payment pages go through `/embedded/payment_checker` for status checking and redirection:

```
https://your-domain.com/embedded/payment_checker?paymentId=pay_xxx&env=prod
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

## Parent Window Integration Code

### 1. Basic Message Listening

```javascript
// Listen for messages from iframe
window.addEventListener('message', function(event) {
  console.log('Received message from iframe:', event.data);
  
  // Security verification: verify message origin (recommend specifying exact domain in production)
  // if (event.origin !== 'https://your-payment-domain.com') return;
  
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
```

### 2. Handler Function Examples

```javascript
function handlePaymentSuccess(paymentId, invoiceId) {
  console.log('Payment successful:', { paymentId, invoiceId });
  
  // Close payment modal
  closePaymentModal();
  
  // Show success message
  showSuccessMessage('Payment completed successfully!');
  
  // Update page status
  updatePaymentStatus('success');
  
  // Optional: redirect to success page
  // window.location.href = '/payment/success';
}

function handlePaymentFailed(paymentId, invoiceId) {
  console.log('Payment failed:', { paymentId, invoiceId });
  
  // Close payment modal
  closePaymentModal();
  
  // Show error message
  showErrorMessage('Payment failed. Please try again.');
  
  // Update page status
  updatePaymentStatus('failed');
}

function handlePaymentCancelled(paymentId, invoiceId) {
  console.log('Payment cancelled:', { paymentId, invoiceId });
  
  // Close payment modal
  closePaymentModal();
  
  // Show info message
  showInfoMessage('Payment was cancelled.');
  
  // Update page status
  updatePaymentStatus('cancelled');
}
```

### 3. Complete iframe Integration Example

```html
<!DOCTYPE html>
<html>
<head>
  <title>Payment Integration Example</title>
  <style>
    .modal {
      display: none;
      position: fixed;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background: rgba(0,0,0,0.5);
      z-index: 1000;
    }
    
    .modal-content {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      background: white;
      border-radius: 8px;
      overflow: hidden;
      width: 90%;
      max-width: 800px;
      height: 80%;
      max-height: 600px;
    }
    
    .modal-header {
      padding: 15px 20px;
      background: #f8f9fa;
      border-bottom: 1px solid #dee2e6;
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
    
    .modal-body {
      height: calc(100% - 60px);
    }
    
    .modal-body iframe {
      width: 100%;
      height: 100%;
      border: none;
    }
    
    .close-btn {
      background: none;
      border: none;
      font-size: 24px;
      cursor: pointer;
    }
  </style>
</head>
<body>
  <button onclick="openPaymentModal()">Open Payment</button>
  
  <div id="paymentModal" class="modal">
    <div class="modal-content">
      <div class="modal-header">
        <h3>Payment</h3>
        <button class="close-btn" onclick="closePaymentModal()">&times;</button>
      </div>
      <div class="modal-body">
        <iframe id="paymentIframe" src="about:blank"></iframe>
      </div>
    </div>
  </div>

  <script>
    // Message listening
    window.addEventListener('message', function(event) {
      console.log('Received message:', event.data);
      
      // Security verification (specify exact domain in production)
      // if (event.origin !== 'https://your-payment-domain.com') return;
      
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
    
    function openPaymentModal() {
      const paymentId = 'pay_xxx'; // Replace with actual payment ID
      const paymentUrl = `https://your-domain.com/embedded/payment_checker?paymentId=${paymentId}&env=prod`;
      
      document.getElementById('paymentIframe').src = paymentUrl;
      document.getElementById('paymentModal').style.display = 'block';
    }
    
    function closePaymentModal() {
      document.getElementById('paymentModal').style.display = 'none';
      document.getElementById('paymentIframe').src = 'about:blank';
    }
    
    function handlePaymentSuccess(paymentId, invoiceId) {
      alert(`Payment successful! Payment ID: ${paymentId}`);
      closePaymentModal();
      // Add your business logic here
    }
    
    function handlePaymentFailed(paymentId, invoiceId) {
      alert(`Payment failed! Payment ID: ${paymentId}`);
      closePaymentModal();
      // Add your business logic here
    }
    
    function handlePaymentCancelled(paymentId, invoiceId) {
      alert(`Payment cancelled! Payment ID: ${paymentId}`);
      closePaymentModal();
      // Add your business logic here
    }
  </script>
</body>
</html>
```

## Security Considerations

### 1. Message Origin Verification

```javascript
window.addEventListener('message', function(event) {
  // Must verify message origin in production
  if (event.origin !== 'https://your-payment-domain.com') {
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

## Testing Tools

We provide `iframe-test.html` testing tool to help you test iframe integration:

1. Open `https://your-domain.com/iframe-test.html`
2. Enter payment ID and base URL
3. Click "Open Payment Modal" to test iframe integration
4. Observe console messages and page behavior

## FAQ

### Q: Why am I not receiving messages?
A: Check the following:
- Ensure iframe loads `/embedded/payment_checker` page
- Check browser console for error messages
- Verify message listener is properly set up
- Validate message origin domain is correct

### Q: How to handle payment timeout?
A: You can set timeout mechanism in parent window:

```javascript
let paymentTimeout;

function openPaymentModal() {
  // Set 30 minutes timeout
  paymentTimeout = setTimeout(() => {
    closePaymentModal();
    showErrorMessage('Payment timeout. Please try again.');
  }, 30 * 60 * 1000);
  
  // Open payment modal...
}

function closePaymentModal() {
  if (paymentTimeout) {
    clearTimeout(paymentTimeout);
    paymentTimeout = null;
  }
  // Close modal...
}
```

### Q: How to customize iframe styles?
A: You can control iframe container styles through CSS:

```css
.payment-iframe-container {
  width: 100%;
  height: 600px;
  border: 1px solid #ddd;
  border-radius: 8px;
  overflow: hidden;
}

.payment-iframe-container iframe {
  width: 100%;
  height: 100%;
  border: none;
}
```

## Technical Support

For issues, please contact technical support or refer to related documentation.

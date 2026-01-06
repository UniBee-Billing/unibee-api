import React, { useState, useEffect, useCallback } from "react";
import { useSearchParams } from "react-router-dom";
import paymentService from '../services/paymentService';

// Loading spinner component
const LoadingSpinner = () => (
  <div style={{
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    padding: '40px 20px'
  }}>
    <div style={{
      width: '50px',
      height: '50px',
      border: '4px solid #f3f3f3',
      borderTop: '4px solid #FFD400',
      borderRadius: '50%',
      animation: 'spin 1s linear infinite',
      marginBottom: '20px'
    }}></div>
    <p style={{
      fontSize: '16px',
      color: '#666',
      margin: 0,
      textAlign: 'center'
    }}>Checking payment status...</p>
  </div>
);

// Success icon component
const SuccessIcon = () => (
  <div style={{
    width: '80px',
    height: '80px',
    borderRadius: '50%',
    backgroundColor: '#4caf50',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    margin: '0 auto 20px',
    animation: 'scaleIn 0.5s ease-out'
  }}>
    <svg width="40" height="40" viewBox="0 0 24 24" fill="white">
      <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
    </svg>
  </div>
);

// Error icon component
const ErrorIcon = () => (
  <div style={{
    width: '80px',
    height: '80px',
    borderRadius: '50%',
    backgroundColor: '#f44336',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    margin: '0 auto 20px',
    animation: 'scaleIn 0.5s ease-out'
  }}>
    <svg width="40" height="40" viewBox="0 0 24 24" fill="white">
      <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
    </svg>
  </div>
);

// Return to Merchant button component
const ReturnToMerchantButton = ({ cancelUrl }) => {
  const handleClick = () => {
    if (cancelUrl) {
      window.location.href = cancelUrl;
    }
  };

  return (
    <button
      onClick={handleClick}
      disabled={!cancelUrl}
      style={{
        backgroundColor: cancelUrl ? '#FFD400' : '#ccc',
        color: 'white',
        border: 'none',
        borderRadius: '8px',
        padding: '12px 24px',
        fontSize: '16px',
        fontWeight: '600',
        cursor: cancelUrl ? 'pointer' : 'not-allowed',
        transition: 'background-color 0.2s ease',
        marginTop: '16px'
      }}
      onMouseOver={(e) => {
        if (cancelUrl) {
          e.target.style.backgroundColor = '#E6C200';
        }
      }}
      onMouseOut={(e) => {
        if (cancelUrl) {
          e.target.style.backgroundColor = '#FFD400';
        }
      }}
    >
      Return to Merchant
    </button>
  );
};

const PaymentChecker = () => {
  const [searchParams] = useSearchParams();
  const paymentId = searchParams.get('paymentId');
  const envParam = searchParams.get('env');
  const env = envParam || 'prod'; // Default to prod if not specified
  
  // Debug env parameter
  console.log('PaymentChecker - envParam:', envParam);
  console.log('PaymentChecker - env (final):', env);
  
  const [, setPayment] = useState(null);

  // Set page title
  useEffect(() => {
    document.title = 'Payment Status Checker - Payment Checkout';
  }, []);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [status, setStatus] = useState('checking'); // checking, success, failed, cancelled
  const [cancelUrl, setCancelUrl] = useState(null);

  // Fetch payment details
  const fetchPaymentDetails = useCallback(async () => {
    try {
      const data = await paymentService.getPaymentDetails(paymentId);
      setPayment(data);
      setError(null);
      return data;
    } catch (err) {
      console.error('Error fetching payment details:', err);
      setError('Payment not found');
      setLoading(false);
      return null;
    }
  }, [paymentId]);

  // Check payment status and handle accordingly
  const checkPaymentStatus = useCallback(async () => {
    const data = await fetchPaymentDetails();
    if (!data || !data.payment) {
      setError('Payment not found');
      setLoading(false);
      return;
    }

    const paymentStatus = data.payment.status;
    console.log('Payment status:', paymentStatus);

    switch (paymentStatus) {
      case 20: // Success
        setStatus('success');
        setLoading(false);
        // Get return URL from data structure
        const returnUrl = data?.returnUrl;
        console.log('Payment success - data:', data);
        console.log('Payment success - returnUrl:', returnUrl);
        console.log('Payment success - env:', env);
        // Send message to parent window (for iframe communication)
        try {
          window.parent.postMessage({
            type: 'UniBee_PaymentSuccess',
            paymentId: paymentId,
            invoiceId: data?.payment?.invoiceId || ''
          }, '*');
          console.log('Payment success - message sent to parent window');
        } catch (error) {
          console.log('Payment success - failed to send message to parent:', error);
        }

        if (returnUrl) {
          // Different redirect timing based on environment
          const redirectDelay = env === 'prod' ? 500 : 2000; // 0.5s for prod, 2s for non-prod
          console.log('Payment success - env check:', env === 'prod');
          console.log('Payment success - redirect delay:', redirectDelay);
          setTimeout(() => {
            console.log('Payment success - redirecting to:', returnUrl);
            window.location.href = returnUrl;
          }, redirectDelay);
        } else {
          console.log('Payment success - no returnUrl found');
        }
        break;
      case 30: // Failed
        setStatus('failed');
        setLoading(false);
        // Send message to parent window (for iframe communication)
        try {
          window.parent.postMessage({
            type: 'UniBee_PaymentFailed',
            paymentId: paymentId,
            invoiceId: data?.payment?.invoiceId || ''
          }, '*');
          console.log('Payment failed - message sent to parent window');
        } catch (error) {
          console.log('Payment failed - failed to send message to parent:', error);
        }
        // Get cancel URL from data structure
        const cancelUrlForFailed = data?.cancelUrl;
        console.log('Payment failed - cancelUrl:', cancelUrlForFailed);
        setCancelUrl(cancelUrlForFailed);
        break;
      case 40: // Cancelled
        setStatus('cancelled');
        setLoading(false);
        // Send message to parent window (for iframe communication)
        try {
          window.parent.postMessage({
            type: 'UniBee_PaymentCancelled',
            paymentId: paymentId,
            invoiceId: data?.payment?.invoiceId || ''
          }, '*');
          console.log('Payment cancelled - message sent to parent window');
        } catch (error) {
          console.log('Payment cancelled - failed to send message to parent:', error);
        }
        // Get cancel URL from data structure
        const cancelUrlForCancel = data?.cancelUrl;
        console.log('Payment cancelled - cancelUrl:', cancelUrlForCancel);
        setCancelUrl(cancelUrlForCancel);
        break;
      case 10: // Pending/Created
      default:
        // Continue polling
        setStatus('checking');
        break;
    }
  }, [fetchPaymentDetails, env, paymentId]);

  // Initial load and polling
  useEffect(() => {
    if (!paymentId) {
      setError('Payment ID is required');
      setLoading(false);
      return;
    }

    // Initial check
    checkPaymentStatus();

    // Set up polling every 1.2 seconds
    const interval = setInterval(() => {
      if (status === 'checking') {
        checkPaymentStatus();
      }
    }, 1200);

    // Cleanup interval on unmount
    return () => clearInterval(interval);
  }, [paymentId, checkPaymentStatus, status]);

  // Render error state
  if (error) {
    return (
      <div className="payment-checker-container" style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        backgroundColor: '#ffffff',
        padding: '20px'
      }}>
        <div className="payment-checker-card" style={{
          backgroundColor: 'white',
          borderRadius: '12px',
          padding: '40px',
          boxShadow: '0 4px 20px rgba(0, 0, 0, 0.1)',
          textAlign: 'center',
          maxWidth: '500px',
          width: '100%'
        }}>
          <ErrorIcon />
          <h2 className="payment-checker-title" style={{
            color: '#f44336',
            margin: '0 0 16px 0',
            fontSize: '24px',
            fontWeight: '600'
          }}>Payment Error</h2>
          <p className="payment-checker-text" style={{
            color: '#666',
            margin: '0 0 24px 0',
            fontSize: '16px',
            lineHeight: '1.5'
          }}>{error}</p>
        </div>
      </div>
    );
  }

  // Render loading state
  if (loading || status === 'checking') {
    return (
      <div className="payment-checker-container" style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        backgroundColor: '#ffffff',
        padding: '20px'
      }}>
        <div className="payment-checker-card" style={{
          backgroundColor: 'white',
          borderRadius: '12px',
          padding: '40px',
          boxShadow: '0 4px 20px rgba(0, 0, 0, 0.1)',
          textAlign: 'center',
          maxWidth: '500px',
          width: '100%'
        }}>
          <LoadingSpinner />
        </div>
      </div>
    );
  }

  // Render success state
  if (status === 'success') {
    return (
      <div className="payment-checker-container" style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        backgroundColor: '#ffffff',
        padding: '20px'
      }}>
        <div className="payment-checker-card" style={{
          backgroundColor: 'white',
          borderRadius: '12px',
          padding: '40px',
          boxShadow: '0 4px 20px rgba(0, 0, 0, 0.1)',
          textAlign: 'center',
          maxWidth: '500px',
          width: '100%'
        }}>
          <SuccessIcon />
          <h2 className="payment-checker-title" style={{
            color: '#4caf50',
            margin: '0 0 16px 0',
            fontSize: '24px',
            fontWeight: '600'
          }}>Payment Successful!</h2>
          <p className="payment-checker-text" style={{
            color: '#666',
            margin: '0 0 24px 0',
            fontSize: '16px',
            lineHeight: '1.5'
          }}>Your payment has been processed successfully. You will be redirected shortly...</p>
        </div>
      </div>
    );
  }

  // Render failed/cancelled state
  if (status === 'failed' || status === 'cancelled') {
    return (
      <div className="payment-checker-container" style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        backgroundColor: '#ffffff',
        padding: '20px'
      }}>
        <div className="payment-checker-card" style={{
          backgroundColor: 'white',
          borderRadius: '12px',
          padding: '40px',
          boxShadow: '0 4px 20px rgba(0, 0, 0, 0.1)',
          textAlign: 'center',
          maxWidth: '500px',
          width: '100%'
        }}>
          <ErrorIcon />
          <h2 className="payment-checker-title" style={{
            color: '#f44336',
            margin: '0 0 16px 0',
            fontSize: '24px',
            fontWeight: '600'
          }}>
            {status === 'failed' ? 'Payment Failed' : 'Payment Cancelled'}
          </h2>
          <p className="payment-checker-text" style={{
            color: '#666',
            margin: '0 0 24px 0',
            fontSize: '16px',
            lineHeight: '1.5'
          }}>
            {status === 'failed' 
              ? 'Your payment could not be processed. Please try again or contact support.'
              : 'Your payment was cancelled. You can try again if needed.'
            }
          </p>
          <ReturnToMerchantButton cancelUrl={cancelUrl} />
        </div>
      </div>
    );
  }

  return null;
};

export default PaymentChecker;

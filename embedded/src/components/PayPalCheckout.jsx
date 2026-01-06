import React, { useState, useEffect, useRef, useCallback } from "react";
import { useSearchParams } from "react-router-dom";
import paymentService from '../services/paymentService';

// PayPal-specific data extraction functions
const extractPayPalOrderId = (paymentData) => {
  const action = paymentService.getActionData(paymentData);
  
  if (action.paypalOrderID) {
    return action.paypalOrderID;
  }
  
  throw new Error('No PayPal Order ID found in payment action data');
};

const getPayPalClientId = (paymentData) => {
  try {
    const action = paymentService.getActionData(paymentData);
    console.log('PayPal action data:', action);
    console.log('PayPal Client ID:', action.paypalClientId);
    return action.paypalClientId || null;
  } catch (error) {
    console.error('Error getting PayPal Client ID:', error);
    return null;
  }
};

const getPaymentCurrency = (paymentData) => {
  try {
    // Get currency from payment data and convert to uppercase
    const currency = paymentData.payment.currency || 'USD';
    return currency.toUpperCase();
  } catch (error) {
    console.error('Error getting payment currency:', error);
    return 'USD';
  }
};

const getPayPalReturnUrl = (paymentData) => {
  try {
    const action = paymentService.getActionData(paymentData);
    return action.paypalReturnUrl || '';
  } catch (error) {
    console.error('Error getting PayPal return URL:', error);
    return '';
  }
};

const getPayPalCancelUrl = (paymentData) => {
  try {
    const action = paymentService.getActionData(paymentData);
    return action.paypalCancelUrl || '';
  } catch (error) {
    console.error('Error getting PayPal cancel URL:', error);
    return '';
  }
};

const PayPalCheckout = () => {
  const [searchParams] = useSearchParams();
  const paymentId = searchParams.get('paymentId');
  const [paymentData, setPaymentData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [paypalLoaded, setPaypalLoaded] = useState(false);
  const paypalRef = useRef();

  // Set page title
  useEffect(() => {
    document.title = 'PayPal Payment - Payment Checkout';
  }, []);

  const renderPayPalButtons = useCallback(() => {
    if (!window.paypal || !paymentData) return;

    // Clear previous buttons
    if (paypalRef.current) {
      paypalRef.current.innerHTML = '';
    }

    const paypalOrderId = extractPayPalOrderId(paymentData);
    if (!paypalOrderId) {
      setError('PayPal Order ID not found in payment data');
      return;
    }

    window.paypal.Buttons({
      createOrder: function(data, actions) {
        // Use the order ID from backend
        return paypalOrderId;
      },
      onApprove: function(data, actions) {
        // Handle successful payment
        return actions.order.capture().then(function(details) {
          console.log('Payment completed:', details);
          
          // Get return URL and redirect
          const returnUrl = getPayPalReturnUrl(paymentData);
          console.log('PayPal return URL:', returnUrl);
          
          if (returnUrl) {
            console.log('Redirecting to return URL:', returnUrl);
            window.location.href = returnUrl;
          } else {
            console.log('No return URL found, showing success message');
            setError('Payment completed successfully! Please close this window.');
          }
        });
      },
      onError: function(err) {
        console.error('PayPal error:', err);
        
        // Get cancel URL and redirect
        const cancelUrl = getPayPalCancelUrl(paymentData);
        console.log('PayPal cancel URL:', cancelUrl);
        
        if (cancelUrl) {
          console.log('Redirecting to cancel URL:', cancelUrl);
          window.location.href = cancelUrl;
        } else {
          setError('Payment failed: ' + err.message);
        }
      },
      onCancel: function(data) {
        console.log('Payment cancelled:', data);
        
        // Get cancel URL and redirect
        const cancelUrl = getPayPalCancelUrl(paymentData);
        console.log('PayPal cancel URL:', cancelUrl);
        
        if (cancelUrl) {
          console.log('Redirecting to cancel URL:', cancelUrl);
          window.location.href = cancelUrl;
        } else {
          setError('Payment was cancelled');
        }
      }
    }).render(paypalRef.current);
  }, [paymentData]);

  const loadPayPalSDK = useCallback((clientId, currency) => {
    // Remove existing PayPal script
    const existingScript = document.querySelector('script[src*="paypal.com/sdk"]');
    if (existingScript) {
      document.body.removeChild(existingScript);
    }

    // Load new PayPal script with correct client ID and currency
    const script = document.createElement('script');
    script.src = `https://www.paypal.com/sdk/js?client-id=${clientId}&currency=${currency}`;
    script.async = true;
    script.onload = () => {
      console.log('PayPal SDK loaded with client ID:', clientId, 'and currency:', currency);
      setPaypalLoaded(true);
    };
    script.onerror = (error) => {
      console.error('Failed to load PayPal SDK:', error);
      setError('Failed to load PayPal SDK');
      setPaypalLoaded(false);
    };
    document.body.appendChild(script);
  }, []);

  // Load PayPal SDK only after payment data is fetched
  useEffect(() => {
    // Don't load SDK here, wait for payment data
    setPaypalLoaded(false);
  }, []);

  // Fetch payment details from backend
  useEffect(() => {
    if (!paymentId) {
      setError('Payment ID is required');
      setLoading(false);
      return;
    }

    const fetchPaymentDetails = async () => {
      try {
        const data = await paymentService.getPaymentDetails(paymentId);
        setPaymentData(data);

        // Get PayPal client ID from payment data
        const paypalClientId = getPayPalClientId(data);
        if (!paypalClientId) {
          setError('PayPal Client ID not found in payment data');
          return;
        }

        // Get payment currency
        const currency = getPaymentCurrency(data);
        console.log('Payment currency:', currency);

        // Load PayPal SDK with correct client ID and currency
        loadPayPalSDK(paypalClientId, currency);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchPaymentDetails();
  }, [paymentId, loadPayPalSDK]);

  // Re-render buttons when payment data or SDK loading status changes
  useEffect(() => {
    if (paypalLoaded && paymentData && window.paypal) {
      console.log('Rendering PayPal buttons...');
      renderPayPalButtons();
    }
  }, [paypalLoaded, paymentData, renderPayPalButtons]);

  if (loading) {
    return (
      <div style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        padding: '16px',
        backgroundColor: '#ffffff'
      }}>
        <div style={{
          width: '100%',
          maxWidth: '540px',
          margin: '0 auto',
          backgroundColor: 'white',
          borderRadius: '12px',
          boxShadow: '0 4px 12px rgba(0, 0, 0, 0.1)',
          padding: '32px',
          textAlign: 'center'
        }}>
          <div style={{ marginBottom: '16px' }}>
            <div style={{
              width: '40px',
              height: '40px',
              border: '4px solid #f3f3f3',
              borderTop: '4px solid #FFD400',
              borderRadius: '50%',
              animation: 'spin 1s linear infinite',
              margin: '0 auto 16px'
            }}></div>
            <h2 style={{ 
              margin: '0 0 8px 0', 
              fontSize: '24px', 
              fontWeight: '700', 
              color: '#1a1a1a' 
            }}>
              Loading Payment
            </h2>
            <p style={{ 
              margin: '0', 
              color: '#666', 
              fontSize: '16px' 
            }}>
              Please wait while we prepare your payment...
            </p>
          </div>
        </div>
        <style>
          {`
            @keyframes spin {
              0% { transform: rotate(0deg); }
              100% { transform: rotate(360deg); }
            }
          `}
        </style>
      </div>
    );
  }

  if (error) {
    return (
      <div style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        padding: '16px',
        backgroundColor: '#ffffff'
      }}>
        <div style={{
          width: '100%',
          maxWidth: '540px',
          margin: '0 auto',
          backgroundColor: 'white',
          borderRadius: '12px',
          boxShadow: '0 4px 12px rgba(0, 0, 0, 0.1)',
          padding: '32px',
          textAlign: 'center'
        }}>
          <div style={{ marginBottom: '24px' }}>
            <div style={{
              width: '48px',
              height: '48px',
              backgroundColor: '#fee2e2',
              borderRadius: '50%',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              margin: '0 auto 16px',
              fontSize: '24px'
            }}>
              ❌
            </div>
            <h2 style={{ 
              margin: '0 0 8px 0', 
              fontSize: '24px', 
              fontWeight: '700', 
              color: '#dc2626' 
            }}>
              Payment Error
            </h2>
            <p style={{ 
              margin: '0 0 24px 0', 
              color: '#666', 
              fontSize: '16px' 
            }}>
              {error}
            </p>
            <button 
              onClick={() => window.location.reload()}
              style={{
                backgroundColor: '#FFD400',
                color: 'white',
                border: 'none',
                borderRadius: '8px',
                padding: '12px 24px',
                fontSize: '16px',
                fontWeight: '600',
                cursor: 'pointer',
                transition: 'background-color 0.2s'
              }}
              onMouseOver={(e) => e.target.style.backgroundColor = '#E6C200'}
              onMouseOut={(e) => e.target.style.backgroundColor = '#FFD400'}
            >
              Try Again
            </button>
          </div>
        </div>
      </div>
    );
  }

  if (!paymentData) {
    return (
      <div style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        padding: '16px',
        backgroundColor: '#ffffff'
      }}>
        <div style={{
          width: '100%',
          maxWidth: '540px',
          margin: '0 auto',
          backgroundColor: 'white',
          borderRadius: '12px',
          boxShadow: '0 4px 12px rgba(0, 0, 0, 0.1)',
          padding: '32px',
          textAlign: 'center'
        }}>
          <div style={{ marginBottom: '24px' }}>
            <div style={{
              width: '48px',
              height: '48px',
              backgroundColor: '#fef3c7',
              borderRadius: '50%',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              margin: '0 auto 16px',
              fontSize: '24px'
            }}>
              ⚠️
            </div>
            <h2 style={{ 
              margin: '0 0 8px 0', 
              fontSize: '24px', 
              fontWeight: '700', 
              color: '#d97706' 
            }}>
              Payment Not Found
            </h2>
            <p style={{ 
              margin: '0', 
              color: '#666', 
              fontSize: '16px' 
            }}>
              The requested payment could not be found. Please check your payment ID and try again.
            </p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div style={{
      minHeight: '100vh',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      padding: '16px',
      backgroundColor: '#f8f9fa'
    }}>
      <div style={{
        width: '100%',
        maxWidth: '540px',
        margin: '0 auto',
        backgroundColor: 'white',
        borderRadius: '12px',
        boxShadow: '0 4px 12px rgba(0, 0, 0, 0.1)',
        padding: '32px',
        textAlign: 'center'
      }}>
        <div style={{ marginBottom: '24px' }}>
          <h2 style={{ 
            margin: '0 0 8px 0', 
            fontSize: '24px', 
            fontWeight: '700', 
            color: '#1a1a1a' 
          }}>
            Complete Payment
          </h2>
          <p style={{ 
            margin: '0', 
            color: '#666', 
            fontSize: '16px' 
          }}>
            Choose your payment method below
          </p>
        </div>
        <div id="paypal-checkout" ref={paypalRef}></div>
      </div>
    </div>
  );
};

export default PayPalCheckout;
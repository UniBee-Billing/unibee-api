import React, { useCallback, useState, useEffect } from "react";
import { loadStripe } from '@stripe/stripe-js';
import {
  EmbeddedCheckoutProvider,
  EmbeddedCheckout
} from '@stripe/react-stripe-js';
import { useSearchParams } from "react-router-dom";
import paymentService from '../services/paymentService';

// Stripe-specific data extraction functions
const extractStripeClientSecret = (paymentData) => {
  const action = paymentService.getActionData(paymentData);
  
  if (action.stripeClientSecret) {
    return action.stripeClientSecret;
  }
  
  throw new Error('No Stripe client secret found in payment action data');
};

const getStripePublishableKey = (paymentData) => {
  try {
    const action = paymentService.getActionData(paymentData);
    return action.stripeAPIKey || null;
  } catch (error) {
    return null;
  }
};

const StripeCheckout = () => {
  const [searchParams] = useSearchParams();
  const paymentId = searchParams.get('paymentId');
  const [paymentData, setPaymentData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [stripePromise, setStripePromise] = useState(null);

  // Set page title
  useEffect(() => {
    document.title = 'Stripe Payment - Payment Checkout';
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

        // Get Stripe publishable key from payment data
        const stripeKey = getStripePublishableKey(data);
        if (stripeKey) {
          setStripePromise(loadStripe(stripeKey));
        } else {
          setError('Stripe API key not found in payment data');
        }
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchPaymentDetails();
  }, [paymentId]);

  const fetchClientSecret = useCallback(() => {
    if (!paymentData) {
      throw new Error('Payment data not available');
    }

    try {
      return Promise.resolve(extractStripeClientSecret(paymentData));
    } catch (err) {
      throw new Error(`Failed to extract client secret: ${err.message}`);
    }
  }, [paymentData]);

  if (loading) {
    return (
      <div style={{ padding: '20px', textAlign: 'center' }}>
        <div>Loading payment details...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div style={{ padding: '20px', textAlign: 'center', color: 'red' }}>
        <div>Error: {error}</div>
        <button onClick={() => window.location.reload()}>Retry</button>
      </div>
    );
  }

  if (!paymentData) {
    return (
      <div style={{ padding: '20px', textAlign: 'center' }}>
        <div>Payment not found</div>
      </div>
    );
  }

  if (!stripePromise) {
    return (
      <div style={{ padding: '20px', textAlign: 'center' }}>
        <div>Initializing payment...</div>
      </div>
    );
  }

  const options = { fetchClientSecret };

  return (
    <div id="checkout" style={{ padding: '20px' }}>
      <EmbeddedCheckoutProvider
        stripe={stripePromise}
        options={options}
      >
        <EmbeddedCheckout />
      </EmbeddedCheckoutProvider>
    </div>
  );
};

export default StripeCheckout;

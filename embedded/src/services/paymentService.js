// Payment service for interacting with the backend API
class PaymentService {
  constructor() {
    this.baseURL = process.env.REACT_APP_API_URL || '';
  }

  // Get dynamic API base URL based on current page path
  getApiBaseURL() {
    const pathname = window.location.pathname;
    console.log('getApiBaseURL - pathname:', pathname);
    const embeddedIndex = pathname.indexOf('/embedded/');
    console.log('getApiBaseURL - embeddedIndex:', embeddedIndex);
    if (embeddedIndex > 0) {
      const prefix = pathname.substring(0, embeddedIndex);
      console.log('getApiBaseURL - prefix:', prefix);
      return prefix;
    }
    console.log('getApiBaseURL - using baseURL:', this.baseURL);
    return this.baseURL;
  }

  /**
   * Fetch payment details by payment ID
   * @param {string} paymentId - The payment ID
   * @returns {Promise<Object>} Payment details including action data
   */
  async getPaymentDetails(paymentId) {
    try {
      const apiBaseURL = this.getApiBaseURL();
      const response = await fetch(`${apiBaseURL}/system/payment/detail?paymentId=${paymentId}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const responseJson = await response.json();
      if (responseJson.code !== 0) {
        throw new Error(responseJson.message);
      }
      return responseJson.data;
    } catch (error) {
      console.error('Error fetching payment details:', error);
      throw error;
    }
  }

  /**
   * Get action data from payment details
   * @param {Object} paymentData - Payment data from backend
   * @returns {Object} Action data object
   */
  getActionData(paymentData) {
    console.log('getActionData - paymentData:', paymentData);
    console.log('getActionData - paymentData.payment:', paymentData?.payment);
    console.log('getActionData - paymentData.payment.action:', paymentData?.payment?.action);
    
    if (!paymentData || !paymentData.payment || !paymentData.payment.action) {
      console.error('getActionData - Missing action data:', {
        hasPaymentData: !!paymentData,
        hasPayment: !!paymentData?.payment,
        hasAction: !!paymentData?.payment?.action
      });
      throw new Error('Payment action data not available');
    }
    return paymentData.payment.action;
  }
}

const paymentService = new PaymentService();
export default paymentService;
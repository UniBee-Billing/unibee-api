import React, { useState, useEffect, useLayoutEffect, useRef, useCallback } from "react";
import { useSearchParams } from "react-router-dom";
import paymentService from '../services/paymentService';

// Message notification component
const Message = ({ type, children, onClose }) => {
  useEffect(() => {
    const timer = setTimeout(() => {
      onClose();
    }, 3000);
    return () => clearTimeout(timer);
  }, [onClose]);

  return (
    <div style={{
      position: 'fixed',
      top: '20px',
      right: '20px',
      padding: '12px 20px',
      backgroundColor: type === 'success' ? '#28a745' : '#dc3545',
      color: 'white',
      borderRadius: '6px',
      boxShadow: '0 4px 12px rgba(0, 0, 0, 0.15)',
      zIndex: 9999,
      fontSize: '14px',
      fontWeight: '500',
      animation: 'slideIn 0.3s ease-out'
    }}>
      {children}
    </div>
  );
};

// Copy button component (unused for now, kept as backup)
// const CopyButton = ({ text, onCopy, children }) => {
//   const handleClick = () => {
//     navigator.clipboard.writeText(text).then(() => {
//       onCopy();
//     }).catch(err => {
//       console.error('Copy failed:', err);
//     });
//   };

//   return (
//     <button 
//       onClick={handleClick}
//       style={{
//         background: 'none',
//         border: 'none',
//         cursor: 'pointer',
//         padding: '8px',
//         borderRadius: '6px',
//         display: 'flex',
//         alignItems: 'center',
//         justifyContent: 'center',
//         transition: 'background-color 0.2s'
//       }}
//       onMouseOver={(e) => e.target.style.backgroundColor = '#f8f9fa'}
//       onMouseOut={(e) => e.target.style.backgroundColor = 'transparent'}
//       title="Copy"
//     >
//       {children}
//     </button>
//   );
// };

// QR Code component
const QRCode = ({ text, onClose }) => {
  // Simple method to generate QR code - using online service
  const qrCodeUrl = `https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(text)}`;
  
  return (
    <div style={{
      position: 'fixed',
      top: '0',
      left: '0',
      right: '0',
      bottom: '0',
      backgroundColor: 'rgba(0, 0, 0, 0.8)',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      zIndex: 10000
    }} onClick={onClose}>
      <div style={{
        backgroundColor: 'white',
        padding: '20px',
        borderRadius: '12px',
        textAlign: 'center',
        maxWidth: '300px',
        width: '90%'
      }} onClick={(e) => e.stopPropagation()}>
        <div style={{ marginBottom: '16px', fontSize: '16px', fontWeight: '600' }}>
          Scan QR Code
        </div>
        <div style={{
          width: '200px',
          height: '200px',
          margin: '0 auto 16px',
          backgroundColor: '#ffffff',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          borderRadius: '8px',
          overflow: 'hidden'
        }}>
          <img 
            src={qrCodeUrl} 
            alt="QR Code" 
            style={{
              width: '100%',
              height: '100%',
              objectFit: 'contain'
            }}
            onError={(e) => {
              e.target.style.display = 'none';
              e.target.nextSibling.style.display = 'flex';
            }}
          />
          <div style={{
            display: 'none',
            width: '100%',
            height: '100%',
            alignItems: 'center',
            justifyContent: 'center',
            fontSize: '12px',
            color: '#666',
            flexDirection: 'column'
          }}>
            <div>QR Code</div>
            <div style={{ fontSize: '10px', wordBreak: 'break-all', padding: '8px' }}>
              {text}
            </div>
          </div>
        </div>
        <div style={{ marginBottom: '16px', fontSize: '12px', color: '#666', wordBreak: 'break-all' }}>
          {text}
        </div>
        <button 
          onClick={onClose}
          style={{
            padding: '8px 16px',
            backgroundColor: '#6c757d',
            color: 'white',
            border: 'none',
            borderRadius: '6px',
            cursor: 'pointer'
          }}
        >
          Close
        </button>
      </div>
    </div>
  );
};

// Blockonomics-specific data extraction functions for Receive Payments API
const getBlockonomicsAction = (paymentData) => {
  try {
    const action = paymentService.getActionData(paymentData);
    return action;
  } catch (error) {
    throw new Error('No Blockonomics action data found in payment data');
  }
};

const getCryptoCurrency = (paymentData) => {
  try {
    return paymentData.payment.cryptoCurrency || 'BTC';
  } catch (error) {
    return 'BTC';
  }
};

const getCryptoAmount = (paymentData) => {
  try {
    const action = getBlockonomicsAction(paymentData);
    const testnet = getTestnet(action);
    const cryptoCurrency = getCryptoCurrency(paymentData);
    
    // If it's test environment and USDT, return 0 amount
    if (testnet === 1 && cryptoCurrency === 'USDT') {
      return 0;
    }
    
    return paymentData.payment.cryptoAmount || 0;
  } catch (error) {
    return 0;
  }
};

const getPaymentAmount = (paymentData) => {
  try {
    const action = getBlockonomicsAction(paymentData);
    const testnet = getTestnet(action);
    const cryptoCurrency = getCryptoCurrency(paymentData);
    
    // If it's test environment and USDT, return 0 amount
    if (testnet === 1 && cryptoCurrency === 'USDT') {
      return 0;
    }
    
    return paymentData.payment.totalAmount || 0;
  } catch (error) {
    return 0;
  }
};

const getCurrency = (paymentData) => {
  try {
    return paymentData.payment.currency || 'USD';
  } catch (error) {
    return 'USD';
  }
};

const getBlockonomicsAddress = (action) => {
  return action.blockonomicsAddress || '';
};

const getBlockonomicsName = (action) => {
  return action.blockonomicsName || 'Payment';
};

const getBlockonomicsDescription = (action) => {
  return action.blockonomicsDescription || '';
};

const getTestnet = (action) => {
  return action.testnet || 0;
};

const getReturnUrl = (action) => {
  return action.blockonomicsReturnUrl || action.returnUrl || '';
};

const getGatewayPaymentId = (paymentData) => {
  try {
    return paymentData.payment.gatewayPaymentId || '';
  } catch (error) {
    return '';
  }
};

// const getCancelUrl = (action) => {
//   return action.cancelUrl || '';
// };

const BlockonomicsCheckout = () => {
  const [searchParams] = useSearchParams();
  const paymentId = searchParams.get('paymentId');
  const [paymentData, setPaymentData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [cryptoCurrency, setCryptoCurrency] = useState('BTC');
  const [paymentStatus, setPaymentStatus] = useState('pending');
  const [lastError, setLastError] = useState('');
  const [paymentCode, setPaymentCode] = useState('');
  const [paymentHTML, setPaymentHTML] = useState('');

  // Set page title
  useEffect(() => {
    document.title = 'Blockonomics Payment - Payment Checkout';
  }, []);
  const [message, setMessage] = useState(null);
  const [showQRCode, setShowQRCode] = useState(false);
  const [qrCodeText, setQrCodeText] = useState('');
  const [originalCurrency, setOriginalCurrency] = useState('');
  const [originalAmount, setOriginalAmount] = useState(0);
  // const [gatewayPaymentId, setGatewayPaymentId] = useState(''); // Used in HTML template, no state needed
  const blockonomicsRef = useRef();
  // const web3PaymentRef = useRef();

  // Show message notification
  const showMessage = (type, text) => {
    setMessage({ type, text });
  };

  // Close message notification
  const closeMessage = () => {
    setMessage(null);
  };

  // Show QR code
  const showQR = (text) => {
    setQrCodeText(text);
    setShowQRCode(true);
  };

  // Close QR code
  const closeQR = () => {
    setShowQRCode(false);
    setQrCodeText('');
  };

  // Handle copy success
  const handleCopySuccess = useCallback(() => {
    showMessage('success', 'Copied successfully!');
  }, []);

  // Format crypto amount for display
  const formatCryptoAmount = (amount, currency) => {
    if (currency === 'BTC') {
      return (amount / 100000000).toFixed(8); // Convert satoshis to BTC
    } else if (currency === 'USDT') {
      return (amount / 1000000).toFixed(6); // Convert to USDT (6 decimals)
    }
    return amount.toString();
  };

  // Effect to set innerHTML when paymentHTML changes and DOM is ready
  useLayoutEffect(() => {
    console.log('useLayoutEffect triggered - paymentHTML:', paymentHTML);
    console.log('useLayoutEffect triggered - blockonomicsRef.current:', blockonomicsRef.current);
    
    if (paymentHTML && blockonomicsRef.current) {
      console.log('Setting innerHTML to blockonomicsRef.current');
      console.log('paymentHTML:', paymentHTML);
      blockonomicsRef.current.innerHTML = paymentHTML;
      console.log('blockonomicsRef.current.innerHTML after setting:', blockonomicsRef.current.innerHTML);
      
      // Set global function for HTML button calls
      window.showCopyAddress = (address) => {
        navigator.clipboard.writeText(address).then(() => {
          handleCopySuccess();
        }).catch(err => {
          console.error('Copy address failed:', err);
          showMessage('error', 'Copy failed, please copy manually');
        });
      };
      
      window.showCopyAmount = (amount) => {
        navigator.clipboard.writeText(amount).then(() => {
          handleCopySuccess();
        }).catch(err => {
          console.error('Copy amount failed:', err);
          showMessage('error', 'Copy failed, please copy manually');
        });
      };
      
      window.showQRCode = (text) => {
        showQR(text);
      };
    } else {
      console.log('useLayoutEffect - conditions not met:', {
        hasPaymentHTML: !!paymentHTML,
        hasBlockonomicsRef: !!blockonomicsRef.current
      });
      
      // If paymentHTML exists but ref is not ready, try again after a short delay
      if (paymentHTML && !blockonomicsRef.current) {
        console.log('Ref not ready, scheduling retry...');
        setTimeout(() => {
          if (blockonomicsRef.current) {
            console.log('Retry - Setting innerHTML to blockonomicsRef.current');
            blockonomicsRef.current.innerHTML = paymentHTML;
            console.log('Retry - blockonomicsRef.current.innerHTML after setting:', blockonomicsRef.current.innerHTML);
          }
        }, 100);
      }
    }
  }, [paymentHTML, handleCopySuccess]);

  // Format fiat amount for display
  const formatFiatAmount = (amount, currency) => {
    return (amount / 100).toFixed(2); // Convert cents to dollars
  };

  // Initialize BTC payment UI
  const initializeBTCPayment = useCallback((data) => {
    try {
      console.log('initializeBTCPayment - Full data:', data);
      console.log('initializeBTCPayment - Payment data:', data.payment);
      console.log('initializeBTCPayment - Action data:', data.payment?.action);
      
      const action = getBlockonomicsAction(data);
      console.log('initializeBTCPayment - Extracted action:', action);
      
      const address = getBlockonomicsAddress(action);
      const name = getBlockonomicsName(action);
      const description = getBlockonomicsDescription(action);
      const cryptoAmount = getCryptoAmount(data);
      const fiatAmount = getPaymentAmount(data);
      const fiatCurrency = getCurrency(data);
      const testnet = getTestnet(action);
      const gatewayPaymentId = getGatewayPaymentId(data);

      console.log('initializeBTCPayment - Address:', address);
      console.log('initializeBTCPayment - CryptoAmount:', cryptoAmount);
      console.log('initializeBTCPayment - FiatAmount:', fiatAmount);
      console.log('initializeBTCPayment - GatewayPaymentId:', gatewayPaymentId);

      if (!address) {
        throw new Error('Bitcoin address not found in payment data');
      }

      const formattedCryptoAmount = formatCryptoAmount(cryptoAmount, 'BTC');
      const formattedFiatAmount = formatFiatAmount(fiatAmount, fiatCurrency);

      // Create BTC payment UI
      const btcPaymentHTML = `
        <div class="blockonomics-card" style="
          text-align: center; 
          padding: 24px; 
          border: 1px solid #dee2e6; 
          border-radius: 12px; 
          background: white;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
          margin-bottom: 16px;
        ">
          <!-- 产品信息和原始金额左右分栏 -->
          <div style="margin: 0 0 24px 0; padding: 16px; background: #f8f9fa; border-radius: 8px; display: flex; justify-content: space-between; align-items: center;">
            <div style="flex: 1;">
              <h3 style="margin: 0 0 4px 0; font-size: 18px; font-weight: 600; color: #1a1a1a;">${name}</h3>
              <p style="margin: 0; font-size: 14px; color: #666;">${description}</p>
            </div>
            <div style="text-align: right;">
              <div style="font-size: 16px; font-weight: 600; color: #FFD400; margin-bottom: 2px;">
                ${(originalAmount / 100).toFixed(2)} ${originalCurrency}
              </div>
              <div style="font-size: 12px; color: #666;">
                Original Order Amount
              </div>
            </div>
          </div>
          
          <div style="margin: 24px 0; padding: 20px; background: #f8f9fa; border-radius: 8px;">
            <div style="display: flex; align-items: center; justify-content: center; gap: 8px; margin-bottom: 4px;">
              <div class="blockonomics-amount" style="font-size: 24px; font-weight: 700; color: #f7931a;">
                ${formattedCryptoAmount} BTC
              </div>
              <button onclick="window.showCopyAmount('${formattedCryptoAmount} BTC')" style="
                background: #f8f9fa; 
                color: #6c757d; 
                border: 1px solid #dee2e6; 
                border-radius: 4px; 
                cursor: pointer;
                padding: 4px 6px;
                display: flex;
                align-items: center;
                justify-content: center;
                transition: all 0.2s;
                font-size: 12px;
              " onmouseover="this.style.backgroundColor='#e9ecef'; this.style.borderColor='#adb5bd'" onmouseout="this.style.backgroundColor='#f8f9fa'; this.style.borderColor='#dee2e6'" title="Copy BTC Amount">
                <svg viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" width="14" height="14" style="fill: currentColor;">
                  <path d="M832 704H384c-35.2 0-64-28.8-64-64V192c0-35.2 28.8-64 64-64h448c35.2 0 64 28.8 64 64v448c0 35.2-28.8 64-64 64zM384 192v448h448V192H384z"></path>
                  <path d="M611.2 896H188.8c-32 0-60.8-28.8-60.8-60.8V412.8c0-32 28.8-60.8 60.8-60.8H256c19.2 0 32 12.8 32 32s-12.8 32-32 32H188.8L192 835.2l419.2-3.2-3.2-64c0-19.2 12.8-32 32-32s32 12.8 32 32v67.2c0 32-28.8 60.8-60.8 60.8zM736 320h-256c-19.2 0-32-12.8-32-32s12.8-32 32-32h256c19.2 0 32 12.8 32 32s-12.8 32-32 32zM736 448h-256c-19.2 0-32-12.8-32-32s12.8-32 32-32h256c19.2 0 32 12.8 32 32s-12.8 32-32 32z"></path>
                  <path d="M736 576h-256c-19.2 0-32-12.8-32-32s12.8-32 32-32h256c19.2 0 32 12.8 32 32s-12.8 32-32 32z"></path>
                </svg>
              </button>
            </div>
            <div style="text-align: center;">
              <span style="font-size: 16px; color: #666;">
                ≈ ${formattedFiatAmount} ${fiatCurrency}
              </span>
            </div>
          </div>
          
          <div style="margin: 24px 0;">
            <label style="display: block; margin-bottom: 12px; font-weight: 600; font-size: 14px; color: #1a1a1a;">Send Bitcoin to:</label>
            <div style="display: flex; align-items: center; gap: 8px; margin-bottom: 12px;">
              <div class="blockonomics-address" style="
                flex: 1;
                background: #f8f9fa; 
                padding: 16px; 
                border: 1px solid #dee2e6; 
                border-radius: 8px; 
                word-break: break-all; 
                font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace; 
                font-size: 13px;
                line-height: 1.4;
                color: #1a1a1a;
              ">
                ${address}
              </div>
              <button onclick="window.showCopyAddress('${address}')" style="
                background: #f8f9fa; 
                color: #6c757d; 
                border: 1px solid #dee2e6; 
                border-radius: 4px; 
                cursor: pointer;
                padding: 6px 8px;
                display: flex;
                align-items: center;
                justify-content: center;
                transition: all 0.2s;
                font-size: 12px;
              " onmouseover="this.style.backgroundColor='#e9ecef'; this.style.borderColor='#adb5bd'" onmouseout="this.style.backgroundColor='#f8f9fa'; this.style.borderColor='#dee2e6'" title="Copy Address">
                <svg viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" width="14" height="14" style="fill: currentColor;">
                  <path d="M832 704H384c-35.2 0-64-28.8-64-64V192c0-35.2 28.8-64 64-64h448c35.2 0 64 28.8 64 64v448c0 35.2-28.8 64-64 64zM384 192v448h448V192H384z"></path>
                  <path d="M611.2 896H188.8c-32 0-60.8-28.8-60.8-60.8V412.8c0-32 28.8-60.8 60.8-60.8H256c19.2 0 32 12.8 32 32s-12.8 32-32 32H188.8L192 835.2l419.2-3.2-3.2-64c0-19.2 12.8-32 32-32s32 12.8 32 32v67.2c0 32-28.8 60.8-60.8 60.8zM736 320h-256c-19.2 0-32-12.8-32-32s12.8-32 32-32h256c19.2 0 32 12.8 32 32s-12.8 32-32 32zM736 448h-256c-19.2 0-32-12.8-32-32s12.8-32 32-32h256c19.2 0 32 12.8 32 32s-12.8 32-32 32z"></path>
                  <path d="M736 576h-256c-19.2 0-32-12.8-32-32s12.8-32 32-32h256c19.2 0 32 12.8 32 32s-12.8 32-32 32z"></path>
                </svg>
              </button>
              <button onclick="window.showQRCode('${address}')" style="
                background: #f8f9fa; 
                color: #6c757d; 
                border: 1px solid #dee2e6; 
                border-radius: 4px; 
                cursor: pointer;
                padding: 6px 8px;
                display: flex;
                align-items: center;
                justify-content: center;
                transition: all 0.2s;
                font-size: 12px;
                font-family: monospace;
              " onmouseover="this.style.backgroundColor='#e9ecef'; this.style.borderColor='#adb5bd'" onmouseout="this.style.backgroundColor='#f8f9fa'; this.style.borderColor='#dee2e6'" title="Show QR Code">
                <svg viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" width="14" height="14" style="fill: currentColor;">
                  <path d="M85.312 85.312V384H384V85.312H85.312zM0 0h469.248v469.248H0V0z m170.624 170.624h128v128h-128v-128zM0 554.624h469.248v469.248H0V554.624z m85.312 85.312v298.624H384V639.936H85.312z m85.312 85.312h128v128h-128v-128zM554.624 0h469.248v469.248H554.624V0z m85.312 85.312V384h298.624V85.312H639.936z m383.936 682.56H1024v85.376h-298.752V639.936H639.936V1023.872H554.624V554.624h255.936v213.248h128V554.624h85.312v213.248z m-298.624-597.248h128v128h-128v-128z m298.624 853.248h-85.312v-85.312h85.312v85.312z m-213.312 0h-85.312v-85.312h85.312v85.312z"></path>
                </svg>
              </button>
            </div>
          </div>
          
          <!-- Gas Fee Warning -->
          <div style="
            margin: 24px 0; 
            padding: 16px; 
            background: ${testnet === 1 ? '#d1ecf1' : '#fff3cd'}; 
            border: 1px solid ${testnet === 1 ? '#bee5eb' : '#ffeaa7'}; 
            border-radius: 8px;
            font-size: 14px;
            line-height: 1.5;
          ">
            ${testnet === 1 ? 
              `<strong style="color: #0c5460;">🧪 Test Mode:</strong> This is a test payment. No actual BTC transfer is required. The payment will be automatically confirmed for testing purposes.` :
              `<strong style="color: #856404;">⚠️ Important:</strong> Please ensure your wallet has sufficient balance to cover network fees (Gas fees). Send exactly <strong>${formattedCryptoAmount} BTC</strong> to the address above. Payment will be confirmed automatically once received.`
            }
          </div>
          
          <!-- Transactions Button -->
          <div style="margin: 24px 0; text-align: right;">
            <button onclick="window.open('https://www.blockonomics.co/#/search?q=${gatewayPaymentId}', '_blank')" style="
              background: #f8f9fa;
              color: #6c757d;
              border: 1px solid #dee2e6;
              border-radius: 6px;
              padding: 8px 16px;
              font-size: 12px;
              font-weight: 500;
              cursor: pointer;
              transition: all 0.2s;
              display: inline-flex;
              align-items: center;
              gap: 6px;
            " onmouseover="this.style.backgroundColor='#e9ecef'; this.style.borderColor='#adb5bd'" onmouseout="this.style.backgroundColor='#f8f9fa'; this.style.borderColor='#dee2e6'" title="View Transaction Details">
              🔍 Transactions
            </button>
          </div>
        </div>
      `;

      // Store the HTML content for later use
      console.log('Setting paymentHTML state with:', btcPaymentHTML);
      setPaymentHTML(btcPaymentHTML);

    } catch (err) {
      console.error('initializeBTCPayment error:', err);
      throw err; // Re-throw the error so it can be caught by the parent function
    }
  }, [originalAmount, originalCurrency]);

  // Helper function to add wallet links
  const addWalletLinks = useCallback((web3Element) => {
    const walker = document.createTreeWalker(
      web3Element,
      NodeFilter.SHOW_TEXT,
      null,
      false
    );
    
    let node;
    while ((node = walker.nextNode())) {
      if (node.textContent.includes('No Web3 wallet found') || 
          node.textContent.includes('Please install browser web3 wallet') ||
          node.textContent.includes('Metamask') ||
          node.textContent.includes('Phantom')) {
        const parent = node.parentNode;
        if (parent && parent.tagName !== 'A' && !parent.innerHTML.includes('href=')) {
          const newHTML = node.textContent
            .replace(/Metamask/g, '<a href="https://metamask.io/download" target="_blank" rel="noopener noreferrer" style="color: #FFD400; text-decoration: underline; font-weight: 600;">Metamask</a>')
            .replace(/Phantom/g, '<a href="https://phantom.app/download" target="_blank" rel="noopener noreferrer" style="color: #FFD400; text-decoration: underline; font-weight: 600;">Phantom</a>');
          
          if (newHTML !== node.textContent) {
            parent.innerHTML = newHTML;
            break; // Found and processed one, stop
          }
        }
      }
    }
  }, []);

  // Load Web3 USDT Component
  const loadWeb3USDTPayment = useCallback((address, amount, returnUrl, testnet) => {
    const container = document.getElementById('web3-payment-container');
    if (!container) return;

    // 1. First show Loading state
    container.innerHTML = `
      <div class="web3-loading" style="
        display: flex;
        flex-direction: column;
        align-items: center;
        padding: 40px 20px;
        text-align: center;
        background: #f8f9fa;
        border-radius: 8px;
        border: 1px solid #e9ecef;
      ">
        <div class="spinner" style="
          width: 40px;
          height: 40px;
          border: 4px solid #f3f3f3;
          border-top: 4px solid #FFD400;
          border-radius: 50%;
          animation: spin 1s linear infinite;
          margin-bottom: 16px;
        "></div>
        <p style="margin: 0; color: #666; font-size: 16px; font-weight: 500;">Connecting Wallet...</p>
        <p style="margin: 8px 0 0 0; color: #999; font-size: 14px;">Please wait while we initialize the payment interface</p>
      </div>
      <style>
        @keyframes spin {
          0% { transform: rotate(0deg); }
          100% { transform: rotate(360deg); }
        }
      </style>
    `;

    // 2. Load Web3 component script
    const script = document.createElement('script');
    script.src = 'https://blockonomics.co/js/web3-payment.js';
    script.async = true;
    
    script.onload = () => {
      // 3. Create Web3 component
      const web3Payment = document.createElement('web3-payment');
      web3Payment.setAttribute('order_amount', amount);
      web3Payment.setAttribute('receive_address', address);
      web3Payment.setAttribute('redirect_url', returnUrl);
      web3Payment.setAttribute('testnet', testnet);
      
      // 4. Add Web3 component to container (after Loading)
      container.appendChild(web3Payment);
      
      // 5. Check if Web3 component has finished loading
      const checkWeb3Completion = () => {
        const web3Element = container.querySelector('web3-payment');
        if (web3Element && web3Element.innerHTML.trim().length > 0) {
          // Hide Loading state
          const loadingElement = container.querySelector('.web3-loading');
          if (loadingElement) {
            loadingElement.style.display = 'none';
          }
          
          // Check if wallet links need to be added
          if (web3Element.textContent.includes('No Web3 wallet found')) {
            addWalletLinks(web3Element);
          }
          
          return true; // Completed
        }
        return false; // Not completed
      };

      // 6. Check immediately once
      if (!checkWeb3Completion()) {
        // If not completed yet, start listening
        const observer = new MutationObserver(() => {
          if (checkWeb3Completion()) {
            observer.disconnect(); // Stop listening after completion
          }
        });
        
        observer.observe(container, { 
          childList: true, 
          subtree: true, 
          characterData: true 
        });
        
        // Set timeout protection
        setTimeout(() => {
          observer.disconnect();
          checkWeb3Completion(); // Final check
        }, 10000);
      }
    };
    
    script.onerror = () => {
      // Show error message when loading fails
      container.innerHTML = `
        <div style="
          display: flex;
          flex-direction: column;
          align-items: center;
          padding: 40px 20px;
          text-align: center;
          background: #f8d7da;
          border-radius: 8px;
          border: 1px solid #f5c6cb;
          color: #721c24;
        ">
          <div style="font-size: 24px; margin-bottom: 16px;">⚠️</div>
          <p style="margin: 0; font-size: 16px; font-weight: 500;">Failed to load Web3 Component</p>
          <p style="margin: 8px 0 0 0; font-size: 14px;">Please refresh the page and try again</p>
        </div>
      `;
    };
    
    document.body.appendChild(script);
  }, [addWalletLinks]);

  // Initialize USDT payment UI with Web3 component
  const initializeUSDTPayment = useCallback((data) => {
    try {
      const action = getBlockonomicsAction(data);
      const address = getBlockonomicsAddress(action);
      const name = getBlockonomicsName(action);
      const description = getBlockonomicsDescription(action);
      const cryptoAmount = getCryptoAmount(data);
      const fiatAmount = getPaymentAmount(data);
      const fiatCurrency = getCurrency(data);
      const testnet = getTestnet(action);
      const returnUrl = getReturnUrl(action);
      const gatewayPaymentId = getGatewayPaymentId(data);

      if (!address) {
        throw new Error('USDT address not found in payment data');
      }

      const formattedCryptoAmount = formatCryptoAmount(cryptoAmount, 'USDT');
      const formattedFiatAmount = formatFiatAmount(fiatAmount, fiatCurrency);

      // Create USDT payment UI
      const usdtPaymentHTML = `
        <div class="blockonomics-card" style="
          text-align: center; 
          padding: 24px; 
          border: 1px solid #dee2e6; 
          border-radius: 12px; 
          background: white;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
          margin-bottom: 16px;
        ">
          <!-- Product Info and Original Amount Side by Side -->
          <div style="margin: 0 0 24px 0; padding: 16px; background: #f8f9fa; border-radius: 8px; display: flex; justify-content: space-between; align-items: center;">
            <div style="flex: 1;">
              <h3 style="margin: 0 0 4px 0; font-size: 18px; font-weight: 600; color: #1a1a1a;">${name}</h3>
              <p style="margin: 0; font-size: 14px; color: #666;">${description}</p>
            </div>
            <div style="text-align: right;">
              <div style="font-size: 16px; font-weight: 600; color: #FFD400; margin-bottom: 2px;">
                ${(originalAmount / 100).toFixed(2)} ${originalCurrency}
              </div>
              <div style="font-size: 12px; color: #666;">
                Original Order Amount
              </div>
            </div>
          </div>
          
          <div style="margin: 24px 0; padding: 20px; background: #f8f9fa; border-radius: 8px;">
            <div style="display: flex; align-items: center; justify-content: center; gap: 8px; margin-bottom: 4px;">
              <div class="blockonomics-amount" style="font-size: 24px; font-weight: 700; color: #26a17b;">
                ${formattedCryptoAmount} USDT
              </div>
              <button onclick="window.showCopyAmount('${formattedCryptoAmount} USDT')" style="
                background: #f8f9fa; 
                color: #6c757d; 
                border: 1px solid #dee2e6; 
                border-radius: 4px; 
                cursor: pointer;
                padding: 4px 6px;
                display: flex;
                align-items: center;
                justify-content: center;
                transition: all 0.2s;
                font-size: 12px;
              " onmouseover="this.style.backgroundColor='#e9ecef'; this.style.borderColor='#adb5bd'" onmouseout="this.style.backgroundColor='#f8f9fa'; this.style.borderColor='#dee2e6'" title="Copy USDT Amount">
                <svg viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" width="14" height="14" style="fill: currentColor;">
                  <path d="M832 704H384c-35.2 0-64-28.8-64-64V192c0-35.2 28.8-64 64-64h448c35.2 0 64 28.8 64 64v448c0 35.2-28.8 64-64 64zM384 192v448h448V192H384z"></path>
                  <path d="M611.2 896H188.8c-32 0-60.8-28.8-60.8-60.8V412.8c0-32 28.8-60.8 60.8-60.8H256c19.2 0 32 12.8 32 32s-12.8 32-32 32H188.8L192 835.2l419.2-3.2-3.2-64c0-19.2 12.8-32 32-32s32 12.8 32 32v67.2c0 32-28.8 60.8-60.8 60.8zM736 320h-256c-19.2 0-32-12.8-32-32s12.8-32 32-32h256c19.2 0 32 12.8 32 32s-12.8 32-32 32zM736 448h-256c-19.2 0-32-12.8-32-32s12.8-32 32-32h256c19.2 0 32 12.8 32 32s-12.8 32-32 32z"></path>
                  <path d="M736 576h-256c-19.2 0-32-12.8-32-32s12.8-32 32-32h256c19.2 0 32 12.8 32 32s-12.8 32-32 32z"></path>
                </svg>
              </button>
            </div>
            <div style="text-align: center;">
              <span style="font-size: 16px; color: #666;">
                ≈ ${formattedFiatAmount} ${fiatCurrency}
              </span>
            </div>
          </div>
          
          <div style="margin: 24px 0;">
            <label style="display: block; margin-bottom: 12px; font-weight: 600; font-size: 14px; color: #1a1a1a;">USDT Address:</label>
            <div style="display: flex; align-items: center; gap: 8px; margin-bottom: 12px;">
              <div class="blockonomics-address" style="
                flex: 1;
                background: #f8f9fa; 
                padding: 16px; 
                border: 1px solid #dee2e6; 
                border-radius: 8px; 
                word-break: break-all; 
                font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace; 
                font-size: 13px;
                line-height: 1.4;
                color: #1a1a1a;
              ">
                ${address}
              </div>
              <button onclick="window.showCopyAddress('${address}')" style="
                background: #f8f9fa; 
                color: #6c757d; 
                border: 1px solid #dee2e6; 
                border-radius: 4px; 
                cursor: pointer;
                padding: 6px 8px;
                display: flex;
                align-items: center;
                justify-content: center;
                transition: all 0.2s;
                font-size: 12px;
              " onmouseover="this.style.backgroundColor='#e9ecef'; this.style.borderColor='#adb5bd'" onmouseout="this.style.backgroundColor='#f8f9fa'; this.style.borderColor='#dee2e6'" title="Copy Address">
                <svg viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" width="14" height="14" style="fill: currentColor;">
                  <path d="M832 704H384c-35.2 0-64-28.8-64-64V192c0-35.2 28.8-64 64-64h448c35.2 0 64 28.8 64 64v448c0 35.2-28.8 64-64 64zM384 192v448h448V192H384z"></path>
                  <path d="M611.2 896H188.8c-32 0-60.8-28.8-60.8-60.8V412.8c0-32 28.8-60.8 60.8-60.8H256c19.2 0 32 12.8 32 32s-12.8 32-32 32H188.8L192 835.2l419.2-3.2-3.2-64c0-19.2 12.8-32 32-32s32 12.8 32 32v67.2c0 32-28.8 60.8-60.8 60.8zM736 320h-256c-19.2 0-32-12.8-32-32s12.8-32 32-32h256c19.2 0 32 12.8 32 32s-12.8 32-32 32zM736 448h-256c-19.2 0-32-12.8-32-32s12.8-32 32-32h256c19.2 0 32 12.8 32 32s-12.8 32-32 32z"></path>
                  <path d="M736 576h-256c-19.2 0-32-12.8-32-32s12.8-32 32-32h256c19.2 0 32 12.8 32 32s-12.8 32-32 32z"></path>
                </svg>
              </button>
              <button onclick="window.showQRCode('${address}')" style="
                background: #f8f9fa; 
                color: #6c757d; 
                border: 1px solid #dee2e6; 
                border-radius: 4px; 
                cursor: pointer;
                padding: 6px 8px;
                display: flex;
                align-items: center;
                justify-content: center;
                transition: all 0.2s;
                font-size: 12px;
                font-family: monospace;
              " onmouseover="this.style.backgroundColor='#e9ecef'; this.style.borderColor='#adb5bd'" onmouseout="this.style.backgroundColor='#f8f9fa'; this.style.borderColor='#dee2e6'" title="Show QR Code">
                <svg viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" width="14" height="14" style="fill: currentColor;">
                  <path d="M85.312 85.312V384H384V85.312H85.312zM0 0h469.248v469.248H0V0z m170.624 170.624h128v128h-128v-128zM0 554.624h469.248v469.248H0V554.624z m85.312 85.312v298.624H384V639.936H85.312z m85.312 85.312h128v128h-128v-128zM554.624 0h469.248v469.248H554.624V0z m85.312 85.312V384h298.624V85.312H639.936z m383.936 682.56H1024v85.376h-298.752V639.936H639.936V1023.872H554.624V554.624h255.936v213.248h128V554.624h85.312v213.248z m-298.624-597.248h128v128h-128v-128z m298.624 853.248h-85.312v-85.312h85.312v85.312z m-213.312 0h-85.312v-85.312h85.312v85.312z"></path>
                </svg>
              </button>
            </div>
          </div>
          
          <!-- Gas Fee Warning -->
          <div style="
            margin: 24px 0; 
            padding: 16px; 
            background: ${testnet === 1 ? '#d1ecf1' : '#fff3cd'}; 
            border: 1px solid ${testnet === 1 ? '#bee5eb' : '#ffeaa7'}; 
            border-radius: 8px;
            font-size: 14px;
            line-height: 1.5;
          ">
            ${testnet === 1 ? 
              `<strong style="color: #0c5460;">🧪 Test Mode:</strong> This is a test payment. No actual USDT transfer is required. The payment will be automatically confirmed for testing purposes.` :
              `<strong style="color: #856404;">⚠️ Important:</strong> Please ensure your wallet has sufficient ETH balance to cover network fees (Gas fees). Send exactly <strong>${formattedCryptoAmount} USDT</strong> to the address above. Payment will be confirmed automatically once received.`
            }
          </div>
          
          <!-- Transactions Button -->
          <div style="margin: 24px 0; text-align: right;">
            <button onclick="window.open('https://www.blockonomics.co/#/search?q=${gatewayPaymentId}', '_blank')" style="
              background: #f8f9fa;
              color: #6c757d;
              border: 1px solid #dee2e6;
              border-radius: 6px;
              padding: 8px 16px;
              font-size: 12px;
              font-weight: 500;
              cursor: pointer;
              transition: all 0.2s;
              display: inline-flex;
              align-items: center;
              gap: 6px;
            " onmouseover="this.style.backgroundColor='#e9ecef'; this.style.borderColor='#adb5bd'" onmouseout="this.style.backgroundColor='#f8f9fa'; this.style.borderColor='#dee2e6'" title="View Transaction Details">
              🔍 Transactions
            </button>
          </div>
          
          <div id="web3-payment-container" style="margin: 24px 0;">
            <!-- Web3 USDT Component will be loaded here -->
          </div>
        </div>
      `;

      // Store the HTML content for later use
      setPaymentHTML(usdtPaymentHTML);
      
      // Load Web3 USDT Component
      loadWeb3USDTPayment(address, formattedCryptoAmount, returnUrl, testnet);

    } catch (err) {
      console.error('initializeUSDTPayment error:', err);
      throw err; // Re-throw the error so it can be caught by the parent function
    }
  }, [originalAmount, originalCurrency, loadWeb3USDTPayment]);


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

        // Get crypto currency and payment status
        const crypto = getCryptoCurrency(data);
        setCryptoCurrency(crypto);

        // Get payment status and error info
        if (data.payment) {
          setPaymentStatus(data.payment.status || 'pending');
          setLastError(data.payment.lastError || '');
          setPaymentCode(data.payment.paymentCode || '');
          
          // Set original currency and amount
          setOriginalCurrency(data.payment.currency || 'USD');
          setOriginalAmount(data.payment.totalAmount || 0);
          
          // gatewayPaymentId is used directly in HTML template, no state needed
        }

        // Initialize payment UI based on crypto currency
        if (crypto === 'BTC') {
          initializeBTCPayment(data);
        } else if (crypto === 'USDT') {
          initializeUSDTPayment(data);
        } else {
          setError(`Unsupported cryptocurrency: ${crypto}`);
        }

      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchPaymentDetails();
  }, [paymentId, initializeBTCPayment, initializeUSDTPayment]);

  // Poll payment status for BTC payments
  useEffect(() => {
    if (cryptoCurrency !== 'BTC' || !paymentId) return;

    const pollInterval = setInterval(async () => {
      try {
        const data = await paymentService.getPaymentDetails(paymentId);
        if (data.payment) {
          const newStatus = data.payment.status;
          const newLastError = data.payment.lastError || '';
          
          if (newStatus !== paymentStatus) {
            setPaymentStatus(newStatus);
            setLastError(newLastError);
            
            // If payment is successful, show success message and redirect
            if (newStatus === 20) { // PaymentSuccess
              setError('');
              
              // Get return URL and redirect after 1.5 seconds
              try {
                const action = getBlockonomicsAction(data);
                const returnUrl = getReturnUrl(action);
                console.log('Payment success - returnUrl:', returnUrl);
                
                if (returnUrl) {
                  setTimeout(() => {
                    console.log('Payment success - redirecting to:', returnUrl);
                    window.location.href = returnUrl;
                  }, 1500); // 1.5 seconds delay
                } else {
                  console.log('Payment success - no returnUrl found');
                }
              } catch (err) {
                console.error('Error getting return URL for redirect:', err);
              }
            }
          }
        }
      } catch (err) {
        console.error('Error polling payment status:', err);
      }
    }, 5000); // Poll every 5 seconds

    return () => clearInterval(pollInterval);
  }, [cryptoCurrency, paymentId, paymentStatus]);

  // Get status display info
  const getStatusDisplay = () => {
    if (paymentStatus === 20) { // PaymentSuccess
      return { text: 'Payment Confirmed', color: '#28a745', icon: '✅' };
    } else if (paymentStatus === 30) { // PaymentFailed
      return { text: 'Payment Failed', color: '#dc3545', icon: '❌' };
    } else if (paymentStatus === 40) { // PaymentCancelled
      return { text: 'Payment Cancelled', color: '#ffc107', icon: '⚠️' };
    } else {
      return { text: lastError || 'Waiting for Payment', color: '#6c757d', icon: '⏳' };
    }
  };

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
        <button onClick={() => window.location.reload()} style={{ marginTop: '10px', padding: '8px 16px', background: '#007bff', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer' }}>
          Retry
        </button>
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

  const statusDisplay = getStatusDisplay();

  return (
    <>
      <style>
        {`
          @keyframes slideIn {
            from {
              transform: translateX(100%);
              opacity: 0;
            }
            to {
              transform: translateX(0);
              opacity: 1;
            }
          }
          
          @media (max-width: 480px) {
            .blockonomics-container {
              padding: 12px !important;
            }
            .blockonomics-card {
              padding: 20px !important;
            }
            .blockonomics-amount {
              font-size: 20px !important;
            }
            .blockonomics-address {
              font-size: 12px !important;
              padding: 12px !important;
            }
            .blockonomics-button {
              padding: 10px 20px !important;
              font-size: 13px !important;
            }
          }
          @media (min-width: 768px) {
            .blockonomics-container {
              padding: 24px !important;
            }
          }
        `}
      </style>
      
      {/* Message notification */}
      {message && (
        <Message type={message.type} onClose={closeMessage}>
          {message.text}
        </Message>
      )}
      
      {/* QR Code popup */}
      {showQRCode && (
        <QRCode text={qrCodeText} onClose={closeQR} />
      )}
      <div style={{
        minHeight: '100vh',
        backgroundColor: '#ffffff',
        padding: '16px',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center'
      }} className="blockonomics-container">
        <div style={{
          width: '100%',
          maxWidth: '540px', // Adjust max width to balance display effect and interface size
          margin: '0 auto'
        }}>
          <div id="blockonomics-checkout">
            {/* Payment Status Header */}
            <div style={{ 
              marginBottom: '24px', 
              textAlign: 'center',
              padding: '16px',
              backgroundColor: 'white',
              borderRadius: '8px',
              boxShadow: '0 1px 3px rgba(0, 0, 0, 0.1)'
            }}>
              <div style={{ fontSize: '18px', fontWeight: '600', color: statusDisplay.color }}>
                {statusDisplay.icon} {statusDisplay.text}
              </div>
              {paymentCode && (
                <div style={{ fontSize: '14px', color: '#666', marginTop: '8px' }}>
                  Transaction: {paymentCode}
                </div>
              )}
            </div>

            {/* Payment UI */}
            <div ref={blockonomicsRef}></div>

            {/* Footer */}
            <div style={{ 
              marginTop: '24px', 
              textAlign: 'center', 
              fontSize: '14px', 
              color: '#666',
              padding: '16px',
              backgroundColor: 'white',
              borderRadius: '8px',
              boxShadow: '0 1px 3px rgba(0, 0, 0, 0.1)'
            }}>
              <div style={{ fontWeight: '500' }}>Powered by Blockonomics</div>
              <div style={{ marginTop: '8px', fontSize: '12px' }}>
                {cryptoCurrency === 'BTC' ? 'Bitcoin Payment' : 'USDT (ETH ERC-20) Payment'}
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
};

export default BlockonomicsCheckout;
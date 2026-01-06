import React from "react";
import {
  BrowserRouter as Router,
  Route,
  Routes
} from "react-router-dom";
import './App.css';
import StripeCheckout from './components/StripeCheckout';
import PayPalCheckout from './components/PayPalCheckout';
import BlockonomicsCheckout from './components/BlockonomicsCheckout';
import PaymentChecker from './components/PaymentChecker';

// Dynamic basename detection
const getBasename = () => {
  const pathname = window.location.pathname;
  const embeddedIndex = pathname.indexOf('/embedded/');
  if (embeddedIndex > 0) {
    return pathname.substring(0, embeddedIndex);
  }
  return '';
};

const App = () => {
  return (
    <div className="App">
      <Router basename={getBasename()}>
        <Routes>
          <Route path="/embedded/stripe" element={<StripeCheckout />} />
          <Route path="/embedded/paypal" element={<PayPalCheckout />} />
          <Route path="/embedded/blockonomics" element={<BlockonomicsCheckout />} />
          <Route path="/embedded/payment_checker" element={<PaymentChecker />} />
          {/* Future payment gateways can be added here */}
          {/* <Route path="/embedded/payment/alipay" element={<AlipayCheckout />} /> */}
        </Routes>
      </Router>
    </div>
  );
};

export default App;
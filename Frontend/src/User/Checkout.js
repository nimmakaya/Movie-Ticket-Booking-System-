import React, { useState, useEffect } from 'react';
import { PayPalScriptProvider, PayPalButtons } from '@paypal/react-paypal-js';

const Checkout = ({ total, onPaymentSuccess }) => {
  const [show, setShow] = useState(false);
  const [success, setSuccess] = useState(false);
  const [orderID, setOrderID] = useState(false);

  // creates a paypal order
  const createOrder = (data, actions) => {
    return actions.order
      .create({
        purchase_units: [
          {
            description: 'Movie Ticket',
            amount: {
              currency_code: 'USD',
              value: total
            },
          },
        ],
      })
      .then(orderID => {
        setOrderID(orderID);
        return orderID;
      });
  };

  // check Approval
  const onApprove = (data, actions) => {
    return actions.order.capture().then(function (details) {
      setSuccess(true);
      onPaymentSuccess(); 
    });
  };

  useEffect(() => {
    if (success) {
      alert('Payment successful!!');
      console.log('Order successful. Your order id is--', orderID);
     
    }
  }, [success, orderID]);

  return (
    <PayPalScriptProvider options={{ 'client-id': 'AcnnE6ahuq5_dzo7MZkootCEe4_mGA-BJN0YqWNTZXgd5L4j-BrNKAyLwrG0zpxKaE0B3KqnTx3lMt3w' }}>
      <div>
        
              <button className="buy-btn" type="submit" onClick={() => setShow(true)}>
                Book Tickets
              </button>
            </div>
   
        <br></br>
        {show ? (
          <PayPalButtons
            style={{ layout: 'vertical' }}
            createOrder={createOrder}
            onApprove={onApprove}
          />
        ) : null}
     
    </PayPalScriptProvider>
  );
};

export default Checkout;

import React, { useEffect, useState } from 'react';
import { useLocation, Link } from 'react-router-dom';
import axios from 'axios';
import './ConfirmationPage.css';
import qrCodeImage from './icons/QR.png';

const ConfirmationPage = () => {
    const location = useLocation();
    const { bookingDetails, yes } = location.state;
    const [venueName, setVenueName] = useState('');
    const [posterURL, setPosterURL] = useState('');

    const [emailSent, setEmailSent] = useState(false);

    useEffect(() => {
      { // Check if email has not been sent yet
        // Fetch venue details and movie poster URL from the API
        const fetchData = async () => {
            try {
                localStorage.setItem('totalPrice', bookingDetails.booking.totalPrice);
                localStorage.setItem('numTickets', bookingDetails.booking.numSeats);
                const venueResponse = await axios.get(`http://localhost:8080/venues/${bookingDetails.booking.venueId}`);
                setVenueName(venueResponse.data.venue_name);

                const movieResponse = await axios.get(`http://localhost:8080/get-movies/${bookingDetails.booking.movieName}`);
                setPosterURL(movieResponse.data.poster_url);

            } catch (error) {
                console.error('Error fetching data:', error);
            }
        };

        fetchData();
    }
    }, []); // Empty dependency array ensures this effect runs only once after mount

    const sendEmailConfirmation = async () => {
        try {
            const userEmail = localStorage.getItem('username');

            const emailBody = `
            <!DOCTYPE html>
            <html lang="en">
            <head>
                <meta charset="UTF-8">
                <meta name="viewport" content="width=device-width, initial-scale=1.0">
                <title>Booking Confirmation</title>
                <style>
                    body {
                        font-family: Arial, sans-serif;
                        background-color: #f4f4f4;
                        padding: 20px;
                    }
                    .container {
                        max-width: 600px;
                        margin: 0 auto;
                        background-color: #fff;
                        border-radius: 10px;
                        overflow: hidden;
                        box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
                    }
                    .header {
                        background-color: #007bff;
                        color: #fff;
                        text-align: center;
                        padding: 20px 0;
                    }
                    .content {
                        padding: 20px;
                    }
                    .content img {
                        max-width: 100%;
                        height: 200px;
                        margin-bottom: 20px;
                        width:200px;
                    }
                    .footer {
                        background-color: #f4f4f4;
                        padding: 20px;
                        text-align: center;
                    }
                </style>
            </head>
            <body>
                <div class="container">
                    <div class="header">
                        <h1>Booking Confirmation</h1>
                    </div>
                    <div class="content">
                    <h2>Booking Confirmation</h2>
                    <p>Venue: ${venueName}</p>
                    <p>Movie: ${bookingDetails.booking.movieName}</p>
                    <p>Show Time: ${bookingDetails.booking.showTime}</p>
                    <p>Date: ${bookingDetails.booking.date}</p>
                    <p>Total Price: $${bookingDetails.booking.totalPrice}</p>
                    <p>Number of Seats: ${yes ? 'Whole screen' : bookingDetails.booking.numSeats}</p>
                    <img src="${posterURL}" alt="Movie Poster" />
                    <img src="https://posters-bmc.s3.eu-north-1.amazonaws.com/QR.png" alt="QR Code" />
                    </div>
                    <div class="footer">
                        <p>Thank you for booking with us!</p>
                    </div>
                </div>
            </body>
            </html>
            
            `;

            await axios.post('http://localhost:8080/send-emails', {
                from: 'bookmycinemaapp@gmail.com',
                to: [userEmail],
                subject: 'Booking Confirmation',
                body: emailBody,
                contentType: 'text/html', 
            });

            console.log('Email confirmation sent successfully!');
        } catch (error) {
            console.error('Error sending email confirmation:', error);
        }
    };

    useEffect(() => {
      // Call sendEmailConfirmation function when all data is available
      if (venueName && posterURL) {
        if (!localStorage.getItem('emailSent')) {
          // Send email confirmation
    sendEmailConfirmation();
    }
    
    
    // Update emailSent state variable to true
    localStorage.setItem('emailSent', true);
      }
  }, [venueName, posterURL]);

    

    return (
        <div className="parent-container">
            <div className="confirmation-container">
                <h2>Booking Confirmation</h2>
                <div className="ticket-details">
                    <h3>Ticket Details:</h3>
                    <p>Venue: {venueName}</p>
                    <p>Movie: {bookingDetails.booking.movieName}</p>
                    <p>Show Time: {bookingDetails.booking.showTime}</p>
                    <p>Date: {bookingDetails.booking.date}</p>
                    <p>Total Price: ${bookingDetails.booking.totalPrice}</p>
                    <p>Number of Seats: {bookingDetails.booking.numSeats}</p>
                    <p>Seat Numbers: {yes ? 'Whole screen' : bookingDetails.booking.seatNumbers.join(',')}</p>
                    <img src={posterURL} alt="Movie Poster" className='poster-img'/>
                    <img src={qrCodeImage} alt="QR Code" className='qr-code'/>
                </div>

                <div className="button-container">
                    <Link to="/welcome">
                        <button className="main-page-button">Go to Main Page</button>
                    </Link>
                    {!yes && (
                        <Link to="/split-tickets">
                            <button className="main-page-button">Split Tickets</button>
                        </Link>
                    )}
                </div>
            </div>
        </div>
    );
};

export default ConfirmationPage;

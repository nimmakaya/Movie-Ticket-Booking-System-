import React, { useEffect, useState } from 'react';
import Cookies from 'js-cookie';

const WelcomePage = () => {
    const [username, setUsername] = useState('');

    useEffect(() => {
        // Retrieve username from cookie when component mounts
        const storedUsername = Cookies.get('username');
        if (storedUsername) {
            setUsername(storedUsername);
        } else {
            // Handle case where username cookie is not set
            // For example, redirect user to login page
            window.location.href = '/login';
        }
    }, []);

    return (
        <div>
            <h2>Welcome, {username}!</h2>
            <p>This is your personalized welcome message.</p>
        </div>
    );
};

export default WelcomePage;

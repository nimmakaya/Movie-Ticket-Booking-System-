import React, { useEffect, useState } from 'react';
import Cookies from 'js-cookie';
import { Link } from 'react-router-dom';

const AdminWelcomePage = () => {
    const [username, setUsername] = useState('');

    useEffect(() => {
        // Retrieve username from cookie when component mounts
        const storedUsername = Cookies.get('username');
        if (storedUsername) {
            setUsername(storedUsername);
        } else {
            // Handle case where username cookie is not set
            // For example, redirect user to login page
            window.location.href = '/admin-login';
        }
    }, []);

    return (
        <div>
            <h2>Welcome Admin, {username}!</h2>
            <p>This is your personalized welcome message.</p>
            <Link to="/admin/create-city">
                <button>Create City</button>
            </Link>
            <Link to="/admin/create-venue">
                <button>Create Venue</button>
            </Link>
        </div>
    );
};

export default AdminWelcomePage;

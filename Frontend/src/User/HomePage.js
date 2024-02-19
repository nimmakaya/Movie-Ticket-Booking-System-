import React from 'react';
import { Link } from 'react-router-dom';

const HomePage = () => {
    return (
        <div>
            <h1>Welcome to Movie Booking</h1>
            <div>
                <Link to="/register">
                    <button>Register</button>
                </Link>
                <Link to="/login">
                    <button>Login</button>
                </Link>
                <Link to="/admin-login">
                    <button>Login as Admin</button>
                </Link>
            </div>
        </div>
    );
};

export default HomePage;

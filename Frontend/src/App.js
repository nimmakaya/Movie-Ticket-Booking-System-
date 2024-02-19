import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Home from './User/HomePage';
import RegisterForm from './User/RegisterForm';
import LoginForm from './User/LoginForm';
import WelcomePage from './User/WelcomePage';
import AdminWelcomePage from './Admin/WelcomePage';
import AdminLoginForm from './Admin/LoginForm';
import CreateCity from './Admin/CreateCity';
import CreateVenue from './Admin/CreateVenue'

const App = () => {
    return (
      <Router>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/register" element={<RegisterForm />} />
          <Route path="/login" element={<LoginForm />} />

          <Route path="/welcome" element={<WelcomePage />} />


          <Route path="/welcome-admin" element={<AdminWelcomePage />} />
          <Route path="/admin-login" element={<AdminLoginForm />} />
          <Route path="/admin/create-city" element={<CreateCity />} />
          <Route path="/admin/create-venue" element={<CreateVenue />} />
        </Routes>
      </Router>
    );
  };

export default App;

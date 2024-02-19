import React, { useState } from 'react';
import axios from 'axios';

const CreateCity = () => {
    const [cityName, setCityName] = useState('');
    const [message, setMessage] = useState('');

    const handleCreateCity = async () => {
        try {
            const response = await axios.post('http://localhost:8080/cities', {  city_name: cityName });
            setMessage(response.data.message);
        } catch (error) {
            console.log(error.response.data)
            setMessage(error.response.data.error);
        }
    };

    return (
        <div>
            <h2>Create New City</h2>
            <input
                type="text"
                placeholder="Enter city name"
                value={cityName}
                onChange={(e) => setCityName(e.target.value)}
            />
            <button onClick={handleCreateCity}>Create City</button>
            {message && <p>{message}</p>}
        </div>
    );
};

export default CreateCity;

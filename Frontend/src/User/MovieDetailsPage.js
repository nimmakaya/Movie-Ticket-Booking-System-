import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useParams, useNavigate } from 'react-router-dom';
import './MovieDetailsPage.css';

const MovieDetailsPage = ({ bodyClassName }) => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [movie, setMovie] = useState(null);
    const [openedMovie, setOpenedMovie] = useState(null);
    const [venueDetails, setVenueDetails] = useState([]);
    const [selectedDate, setSelectedDate] = useState(""); 
    const [videoId, setVideoId] = useState("");

    useEffect(() => {
        const fetchMovieDetails = async () => {
            try {
                const response = await axios.get(`http://localhost:8080/movies/${id}`);
                setMovie(response.data);
                setVideoId(response.data.trailer_url.split("/").pop().split("?")[0])
                const today = new Date().toISOString().split('T')[0]; // Get today's date in the required format
                setSelectedDate(today);
                fetchOpenedMovie(today);
            } catch (error) {
                console.error('Error fetching movie details:', error);
            }
        };
    
        fetchMovieDetails();
    }, [id]);
    
    const handleDateChange = async (e) => {
        const selectedDate = e.target.value;
        setSelectedDate(selectedDate);
        fetchOpenedMovie(selectedDate);
    };
    
    const fetchOpenedMovie = async (selectedDate) => {
        try {
            const response = await axios.get(`http://localhost:8080/opened_movies?movie_id=${id}&date=${selectedDate}`);
            setOpenedMovie(response.data);
        } catch (error) {
            console.error('Error fetching opened movie:', error);
        }
    };

    const handleShowTimeSelection = (selectedMovie) => {
        navigate(`/seatlayout`, { state: { movie: selectedMovie }});
    };
    
    useEffect(() => {
        const fetchVenueDetails = async () => {
            if (openedMovie && openedMovie.length > 0 && openedMovie[0]?.venues) {
                const venuePromises = openedMovie[0].venues.map(venueId =>
                    axios.get(`http://localhost:8080/venues/${venueId}`)
                );
        
                try {
                    const venueResponses = await Promise.all(venuePromises);
                    const venueDetails = venueResponses.map(response => response.data);
                    setVenueDetails(venueDetails);
                } catch (error) {
                    console.error('Error fetching venue details:', error);
                }
            }
        };        

        fetchVenueDetails();
    }, [openedMovie]);

    if (!movie) {
        return <div>Loading...</div>;
    }

    const releaseDate = new Date(movie.release_date).toDateString();

    const isTimeDisabled = (showTime) => {
        const [hoursStr, minutesStr] = showTime.split(/[:apm]+/);
        let hours = parseInt(hoursStr);
        const minutes = parseInt(minutesStr);
    
        if (showTime.includes('pm') && hours !== 12) {
            hours += 12;
        }
    
       // const selectedDateTime = new Date(selectedDate);
        const [year, month, day] = selectedDate.split('-').map(Number);
        const selectedDateTime = new Date(Date.UTC(year, month - 1, day));
        
        selectedDateTime.setHours(selectedDateTime.getHours() + 4);


         console.log(selectedDateTime)
         selectedDateTime.setHours(hours);
    
        const currentTime = new Date();
    
        return selectedDateTime < currentTime;
    };
    
    

    document.body.className = bodyClassName || '';

    return (
        <div className="welcome-container">
            <h2>Book My Cinema</h2><br />
            <div className="movie-details-container">
                <h3 className="title">{movie.name}</h3>
                <div className="details">
                    <img src={movie.poster_url} alt={movie.name} className="poster" />
                    <p><span className="field">Release Date:</span> {releaseDate}</p>
                    <p><span className="field">Cast:</span> {movie.cast.join(', ')}</p>
                    <p><span className="field">Crew:</span> {movie.crew.join(', ')}</p><br></br>
                    <p>Watch trailer here</p>
                    <iframe
        width="560"
        height="315"
        src={`https://www.youtube.com/embed/${videoId}`}
        title="YouTube video player"
        frameBorder="0"
        allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
        allowFullScreen
    ></iframe><br></br><br></br>
                    <div className="venues">
                        <label htmlFor="datePicker" style={{ marginRight: '10px' }}>Select Date:</label>
                        <input
                            type="date"
                            id="datePicker"
                            value={selectedDate}
                            min={new Date().toISOString().split('T')[0]} 
                            onChange={handleDateChange}
                            style={{ width: '150px' }} 
                        />

                        <h4>Venues:</h4>
                        {openedMovie && openedMovie.length > 0 ? (
                            venueDetails.length > 0 ? (
                                venueDetails.map(venue => (
                                    <div key={venue._id} className="venue">
                                        <p><strong>{venue.venue_name}</strong></p>
                                        <p>{venue.address}</p>
                                        <p><em>Show Times:</em> 
                                            {openedMovie[0].show_times.map(showTime => (
                                                <span key={showTime} className={`${isTimeDisabled(showTime) ? 'disabled' : 'show-time-link'}`} onClick={() => handleShowTimeSelection({ venue_id: venue._id, show_time: showTime, date: selectedDate, movie: movie.name})}>{showTime}</span>
                                            ))}
                                        </p>
                                    </div>
                                ))
                            ) : (
                                <p>No opened venues available for selected date.</p>
                            )
                        ) : (
                            <p>No opened movie available for selected date.</p>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
};

export default MovieDetailsPage;

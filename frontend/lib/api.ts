import axios from 'axios';

export const api = axios.create({
    baseURL: 'http://localhost:8080',
    withCredentials: true, // Essential for sending/receiving session cookies
    headers: {
        'Content-Type': 'application/json',
    },
});
import axios from "axios";
import Cookies from 'js-cookie';

export async function authenticate(email: string, password: string ) {
  try {
    const response = await axios.post('http://localhost:8080/login', {
      email,
      password
    }, {
      withCredentials: true
    });

    const cookie = response.headers['set-cookie'];
    console.log(cookie);

    // Set the cookie using js-cookie
    Cookies.set('first_name', response.data.data.first_name, { path: '/' });

    console.log(response.data);
    console.log(Cookies.get('first_name'));
  } catch (error) {
    console.error('Authentication failed:', error);
  }
}

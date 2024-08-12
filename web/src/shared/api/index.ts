import { API_URL } from "../../constants"
import axios from "axios"

const $api = axios.create({
    baseURL: API_URL,
    withCredentials: true,
})

$api.interceptors.request.use((config) => {
    config.headers.Authorization = `Bearer ${localStorage.getItem("AUTH/ACCESS_TOKEN")}`
    return config
})

$api.interceptors.response.use((config) => {
    return config
}, async (error) => {
    const originalRequest = error.config
    if (error.response.status === 401 && error.config && !error.config._isRetry) {
        originalRequest._isRetry = true
        try {
            console.log("Refreshing token from interceptor")
            const response = await axios.get(`${API_URL}/auth/refreshToken`, { withCredentials: true })
            localStorage.setItem("token", response.data.accessToken)
            return $api.request(originalRequest)
        } catch (e) {
            console.log("User not authorized")
        }
    }
    throw error
})

export default $api
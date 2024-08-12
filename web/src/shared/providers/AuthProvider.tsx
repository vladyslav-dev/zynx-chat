import { ReactNode, useEffect } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import useAuthStore from "../store/auth";

const AuthProvider = ({ children }: { children: ReactNode}) => {
    const location = useLocation();
    const navigate = useNavigate()
    const { isAuthenticated, checkAuth, isAuthenticationVerified, setAuthenticationVerified } = useAuthStore()

    useEffect(() => {
        if (localStorage.getItem('AUTH/ACCESS_TOKEN')) {
            checkAuth()
        } else {
            // User is not authenticated
            setAuthenticationVerified()
        }
    }, [])

    useEffect(() => {
        if (!isAuthenticationVerified) {
            return
        }

        if (isAuthenticated) {
            if (location.pathname === '/login' || location.pathname === '/register') {
                navigate('/')
            }
        } else {
            if (location.pathname !== '/login' && location.pathname !== '/register') {
                navigate('/login')
            }
        }
    }, [isAuthenticated, isAuthenticationVerified])

    return isAuthenticationVerified ? children : <div>Loading...</div>
}

export default AuthProvider
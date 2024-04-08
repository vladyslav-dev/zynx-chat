import { IUser } from "../types/User";
import { createContext, useEffect, useState } from "react";
import { router } from "./RouterProvider";


interface IAuthContext {
    isAuthenticated: boolean;
    setIsAuthenticated: (isAuthenticated: boolean) => void;
    user: IUser | null;
    setUser: (user: IUser) => void;
}

const AuthContext = createContext<IAuthContext>({
    isAuthenticated: false,
    setIsAuthenticated: (isAuthenticated: boolean) => {},
    user: null,
    setUser: (user: IUser) => {}
})

export const AuthContextProvider = ({ children }: { children: React.ReactNode}) => {
    const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false)
    const [user, setUser] = useState<IUser | null>(null)

    useEffect(() => {
        const user = localStorage.getItem('user')

        if (user) {
            setUser(JSON.parse(user))
            setIsAuthenticated(true)

            router.navigate('/dashboard')
        } else {
            setUser(null)
            setIsAuthenticated(false)

            router.navigate('/login')
        }

    }, [isAuthenticated])


    const value = {
        isAuthenticated,
        setIsAuthenticated,
        user,
        setUser
    }

    return (
        <AuthContext.Provider value={value}>
            {children}
        </AuthContext.Provider>
    )
}

export default AuthContextProvider;
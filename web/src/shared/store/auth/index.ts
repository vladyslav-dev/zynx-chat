import { TLoginData, TRegisterData } from '../../api/auth/auth.types';
import * as authService from '../../api/auth/auth.service';
import { create } from "zustand"

type User = {
    id: number
    username: string
    phone: string
    access_token: string
}

interface IAuthStore {
    user: User | null,
    isAuthenticated: boolean,
    isAuthenticationVerified: boolean,
    authError: string | null,
    

    register: (data: TRegisterData) => Promise<void>,
    login: (data: TLoginData) => Promise<void>,
    logout: () => Promise<void>,
    checkAuth: () => Promise<void>
    setAuthenticationVerified: () => void
}

const useAuthStore = create<IAuthStore>((set, get) => ({
    user: null,
    isAuthenticated: false,
    isAuthenticationVerified: false,
    authError: null,

    register: async (data: TRegisterData) => {
        try {
            const response = await authService.register(data)

            if (response.status === 201) {
                get().login({ phone: data.phone, password: data.password })
            }
        } catch (error: any) {
            console.error('Error occured while registering.', error)

            set({ 
                user: null, 
                isAuthenticated: false, 
                authError: error.response.data.message || 'Error occured while registering' 
            })
        }
    },
    login: async (data: TLoginData) => {
        try {
            const response = await authService.login(data)
            const user = response.data

            localStorage.setItem('AUTH/ACCESS_TOKEN', user.access_token)
            set({ user, isAuthenticated: true, authError: null })

            console.log('Successfully logged in', response)
        } catch (error: any) {
            if (error?.response?.status === 409) {
                /* User already logged in */
                console.error('User already logged in')
                console.info('Logging out user')

                get().logout()
                return
            }

            set({ 
                user: null, 
                isAuthenticated: false,
                authError: error.response.data.message || 'Error occured while logging in' 
            })
            console.error('Error occured while logging in.', error)
        }
    },
    logout: async () => {
        try {
            await authService.logout()
        } catch (error) {
            console.error('Error occured while logging out.', error)
        } finally {
            localStorage.removeItem('AUTH/ACCESS_TOKEN')
            set({ user: null, isAuthenticated: false })
        }
    },
    checkAuth: async () => {
        try {
            const response = await authService.checkAuth()
            const user = response.data

            localStorage.setItem('AUTH/ACCESS_TOKEN', user.access_token)
            set({ user, isAuthenticated: true, isAuthenticationVerified: true, authError: null })

        } catch (error: any) {
            console.error('User is not authenticated', error)

            localStorage.removeItem('AUTH/ACCESS_TOKEN')
            set({ user: null, isAuthenticated: false, isAuthenticationVerified: true })
        }
    },
    setAuthenticationVerified: () => {
        set({ isAuthenticationVerified: true })
    }
    
}));

export default useAuthStore
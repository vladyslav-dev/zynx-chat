import axios from "axios"
import $api from "../index"
import { TLoginData, TRegisterData } from "./auth.types"
import { API_URL } from "../../../constants"


export const register = async (data: TRegisterData) => {
    return $api.post("/auth/register", data)
}

export const login = async (data: TLoginData) => {
    return $api.post("/auth/login", data)
}

export const logout = async () => {
    return $api.post("/auth/logout")
}

export const checkAuth = async () => {
    return axios.get(`${API_URL}/auth/refreshToken`, { withCredentials: true })
}
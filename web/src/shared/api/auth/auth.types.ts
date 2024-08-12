import { RequireAtLeastOne } from "../../lib/utilTypes"

export type TAuthData = {
    username: string
    phone: string
    password: string
}

export type TRegisterData = TAuthData

export type TLoginData = RequireAtLeastOne<TAuthData, "username" | "phone">
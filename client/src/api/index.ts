import { API_URL } from "../constants"

export const registerUser = async (data: any) => {
    try {
        const response = await fetch(`${API_URL}/api/register`, {
            headers: {
                'Content-Type': 'application/json'
            },
            method: 'POST',
            body: JSON.stringify(data)
        })

        const responseData = await response.json()

        return responseData

    } catch (err) {
        console.error(err)
    }
}

export const loginUser = async (data: any) => {
    try {
        const response = await fetch(`${API_URL}/api/login`, {
            headers: {
                'Content-Type': 'application/json'
            },
            method: 'POST',
            body: JSON.stringify(data),
            credentials: 'include'
        })

        const responseData = await response.json()

        return responseData

    } catch (err) {
        console.error(err)
    }
}
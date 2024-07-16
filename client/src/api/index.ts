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
        })

        const responseData = await response.json()

        return responseData

    } catch (err) {
        console.error(err)
    }
}

export const getAllUsers = async () => {
    try {
        const response = await fetch(`${API_URL}/api/getAllUsers`, {
            method: 'GET',
        })

        const responseData = await response.json()

        return responseData

    } catch (err) {
        console.error(err)
    }
}

export const getAllGroups = async () => {
    try {
        const response = await fetch(`${API_URL}/api/getAllGroups`, {
            method: 'GET',
        })
        const data = await response.json()

        return data
    } catch (err) {
        console.error(err)
    }
}

export const createGroup = async (data: any) => {
    try {
        const response = await fetch(`${API_URL}/api/createGroup`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        })

        return response

    } catch (err) {
        console.error(err)
    }
}

export const getPrivateMessages = async ({ senderId, recipientId }: { senderId: string, recipientId: string }) => {
    try {
        const response = await fetch(`${API_URL}/api/private-message`, {
            method: "POST",
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ sender_id: senderId, recipient_id: recipientId })
        })

        const data = await response.json()

        return data
    } catch (err) {
        console.error(err)
    }
}

export const getGroupMessages = async (groupId: string) => {
    try {
        const response = await fetch(`${API_URL}/api/group-message`, {
            method: "POST",
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ group_id: groupId})
        })

        const data = await response.json()

        return data
    } catch (err) {
        console.error(err)
    }
}
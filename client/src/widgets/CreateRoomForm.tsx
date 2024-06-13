import { API_URL } from "../constants"
import React from "react"
import { v4 as uuidv4 } from 'uuid'

interface ICreateRoomForm {
    onRoomCreate: () => void
}

const CreateRoomForm = ({
    onRoomCreate
}: ICreateRoomForm) => {

    const onSubmit = async (e: React.FormEvent) => {
        e.preventDefault()

        const formData = new FormData(e.target as HTMLFormElement);

        const { roomName } = Object.fromEntries(formData.entries())

        console.log(roomName)

        if (!roomName) {
            alert('Room name is required');
            return
        }

        const response = await fetch(`${API_URL}/ws/createRoom`, {
            method: 'POST',
            headers: {
                'Content-Type': "application/json"
            },
            credentials: 'include',
            body: JSON.stringify({
                id: uuidv4(),
                name: roomName
            })
        })

        const data = await response.json()
        console.log('onRoomCreate', data)
        if (response.ok) {
            onRoomCreate()
        } else {
            console.error(data)
        }
    }

    return (
        <div>
            <h1>Create Room</h1>
            <form onSubmit={onSubmit}>
                <input type="text" placeholder="Room Name" name="roomName" />
                <button type="submit">Create</button>
            </form>
        </div>
    )
}

export default CreateRoomForm
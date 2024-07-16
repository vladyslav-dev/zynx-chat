import { createGroup } from "../api"
import { API_URL } from "../constants"
import React from "react"
import { v4 as uuidv4 } from 'uuid'

interface ICreateGroupForm {
    onRoomCreate: () => void
}

const CreateGroupForm = ({
    onRoomCreate
}: ICreateGroupForm) => {

    const onSubmit = async (e: React.FormEvent) => {
        e.preventDefault()

        const formData = new FormData(e.target as HTMLFormElement);

        const { groupName } = Object.fromEntries(formData.entries())

        if (!groupName) {
            alert('Room name is required');
            return
        }

        const data = await createGroup({ name: groupName })

        console.log('onRoomCreate', data)
        if (data?.ok) {
            onRoomCreate()
        } else {
            console.error(data)
        }
    }

    return (
        <div>
            <form onSubmit={onSubmit}>
                <input type="text" placeholder="Room Name" name="groupName" />
                <button type="submit">Create</button>
            </form>
        </div>
    )
}

export default CreateGroupForm
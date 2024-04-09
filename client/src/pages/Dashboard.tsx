import CreateRoomForm from "@/widgets/CreateRoomForm"
import { useState } from "react"

const Dashboard = () => {
    const [rooms, setRooms] = useState([] as any[])

    const getAllRooms = () => {

    }

    return (
        <div>
            <h1>Dashboard</h1>

            <CreateRoomForm onRoomCreate={getAllRooms} />

        </div>
    )
}

export default Dashboard
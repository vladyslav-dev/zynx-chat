import { WebsocketContext } from "../providers/WebsocketProvider"
import CreateRoomForm from "../widgets/CreateRoomForm"
import { useContext, useEffect, useState } from "react"
import { API_URL, WEBSOCKET_URL } from "../constants"
import { AuthContext } from "../providers/AuthProvider"
import { router } from "../providers/RouterProvider"

const JoinRoom = () => {
    const [rooms, setRooms] = useState([] as any[])
    const { user } = useContext(AuthContext);
    const { websocketConnection, setWebsocketConnection } = useContext(WebsocketContext)

    const getAllRooms = async () => {
        try {
            const response = await fetch(`${API_URL}/ws/getRooms`, {
                method: 'GET',
            })
            const data = await response.json()

            if (response.ok) {
                setRooms(data)
            } else {
                console.error(data)
            }
        } catch (err) {
            console.error(err)
        }
    }

    const joinRoom = (roomId: string) => {
        const ws = new WebSocket(`${WEBSOCKET_URL}/ws/joinRoom/${roomId}?clientId=${user?.id}&username=${user?.username}`)

        if (ws.OPEN) {
            setWebsocketConnection(ws)


            router.navigate(`/room/${roomId}`)
            // Redirect to room/roomId page
        }
    }

    useEffect(() => {
        getAllRooms()
    }, [])
    console.log('rooms', rooms)
    return (
        <div>
            <h1>Join Room</h1>

            <CreateRoomForm onRoomCreate={getAllRooms} />

            <h3 className="mt-4">Available Rooms:</h3>
            {rooms.map((room: any) => (
                <div key={room.id} className="border border-1 p-2">
                    <p>{room.name}</p>
                    <button onClick={() => joinRoom(room.id)}>Join</button>
                </div>
            ))}


        </div>
    )
}

export default JoinRoom
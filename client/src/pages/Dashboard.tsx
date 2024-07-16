import { WebsocketContext } from "../providers/WebsocketProvider"
import CreateGroupForm from "../widgets/CreateGroupForm"
import { useContext, useEffect, useState } from "react"
import { API_URL, WEBSOCKET_URL } from "../constants"
import { AuthContext } from "../providers/AuthProvider"
import { router } from "../providers/RouterProvider"
import { getAllGroups, getAllUsers } from "../api"
import CurrentUser from "../widgets/CurrentUser"

const Dashboard = () => {
    const [groups, setGroups] = useState([] as any[])
    const [users, setUsers] = useState([] as any[])
    const { user } = useContext(AuthContext);
    

    useEffect(() => {
        if (user) {
            fetchAllData()
        }
    }, [user])

    const fetchAllData = async () => {
        const [users, groups] = await Promise.all([getAllUsers(), getAllGroups()]);
        console.log('users', users)
        console.log('user', user)
        const filteredUsers = users.filter((u: any) => String(u.id) !== String(user?.id));
        setUsers(filteredUsers);
        setGroups(groups);
    }


    const openChat = ({ type, id, name }: { type: string, id: string, name: string }) => {
        router.navigate(`/chat?type=${type}&name=${name}&id=${id}`)
    }

    return (
        <div className="p-4">
            <CurrentUser />

            <h1 className="mb-4 text-4xl">Dashboard</h1>

            <section className="mb-4">
                <h3 className="mb-2">Users</h3>
                {users.map((user: any) => (
                    <div key={user.id} className="border border-1 p-2 flex justify-between">
                        <div>
                            <div className="text-red-400">ID: {user.id}</div>
                            <div>Username: {user.username}</div>
                            <div>Email: {user.email}</div>
                        </div>
                        <button onClick={() => openChat({ type: "private", id: user.id, name: user.username })}>Open chat</button>
                    </div>
                ))}
            </section>

            <section className="mb-4">
                <h3 className="mb-2">Groups</h3>
                <CreateGroupForm onRoomCreate={fetchAllData} />
                {groups.map((group: any) => (
                    <div key={group.id} className="border border-1 p-2 flex justify-between">
                        <div>
                            <div className="text-red-400">ID: {group.id}</div>
                            <div>Name: {group.name}</div>
                        </div>
                        <button onClick={() => openChat({ type: "group", id: group.id, name: group.name })}>Open chat</button>
                    </div>
                ))}
            </section>

            



        </div>
    )
}

export default Dashboard
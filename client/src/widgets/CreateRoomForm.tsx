
interface ICreateRoomForm {
    onRoomCreate: () => void
}

const CreateRoomForm = ({
    onRoomCreate
}: ICreateRoomForm) => {

    const onSubmit = (e: React.FormEvent) => {
        e.preventDefault()

        const formData = new FormData(e.target as HTMLFormElement);

        const { roomName } = Object.fromEntries(formData.entries())

        console.log(roomName)

        onRoomCreate()
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
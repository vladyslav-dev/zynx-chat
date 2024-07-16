import { useLocation } from "react-router-dom";
import { AuthContext } from "../providers/AuthProvider";
import { useContext, useEffect, useRef, useState } from "react";
import { WebsocketContext } from "../providers/WebsocketProvider";
import CurrentUser from "../widgets/CurrentUser";
import { getGroupMessages, getPrivateMessages } from "../api";

type Sender = {
    id: string;
    username: string;
    email: string;
}

type Recipient = {
    id: string;
    username: string;
    email: string;
}

type PrivateMessage = {
    id: number;
    type: 'private';
    sender: Sender;
    recipient: Recipient;
    content: string;
    created_at: string;
};

type GroupMessage = {
    id: number;
    type: 'group';
    sender: Sender;
    group: {
        id: string;
        name: string;
    };
    content: string;
    created_at: string;
};

type Message = PrivateMessage | GroupMessage

const Chat = () => {
    const { user } = useContext(AuthContext);
    const [messages, setMessages] = useState<Message[]>([])
    const [inputMessage, setInputMessage] = useState<string>('')
    const wsRef = useRef<WebSocket | null>(null)

    const location = useLocation();
    const queryParams = new URLSearchParams(location.search);

    
    const type = queryParams.get("type");
    const id = queryParams.get("id");
    const name = queryParams.get("name");
    
    useEffect(() => {
        (async () => {
            if (type === "private") fetchPrivateMessages()
            if (type === "group") fetchGroupMessages()
        })()

    if (!user) return

        const WS_PRIVATE_URL = `ws://localhost:8080/ws/message?type=private&sender_id=${user.id}&recipient_id=${id}`
        const WS_GROUP_URL = `ws://localhost:8080/ws/message?type=group&group_id=${id}`
        wsRef.current = new WebSocket(type === "private" ? WS_PRIVATE_URL : WS_GROUP_URL)

        wsRef.current.onopen = () => {
            console.log('ws opened')
        }

        wsRef.current.onmessage = (event) => {
            console.log('event data', event.data)

            if (type === "private") fetchPrivateMessages()
            if (type === "group") fetchGroupMessages()
        }   

        
    }, [])

    const fetchPrivateMessages = async () => {
        if (!id || !user?.id) return

        const messages = await getPrivateMessages({ senderId: user.id, recipientId: id })
        setMessages(messages)
    }

    const fetchGroupMessages = async () => {
        if (!id) return

        const messages = await getGroupMessages(id)
        setMessages(messages)
    }
    console.log('messages', messages)

    const sendMessage = () => {
        console.log('sendMessage', inputMessage)
        if (!wsRef.current) return

        let message

        if (type === "private") {
            message = {
                "type": "private",
                "recipient_id": parseInt(id!),
                "sender_id": parseInt(user!.id),
                "group_id": 1, // will not be used. Add this to aviod error
                "content": inputMessage
            }
        } else {
            message = {
                "type": "group",
                "recipient_id": 1, // will not be used. Add this to aviod error
                "sender_id": parseInt(user!.id),
                "group_id": parseInt(id!), 
                "content": inputMessage
            }
        }

        console.log('message', message)
        wsRef.current.send(JSON.stringify(message));
    }

    return (
        <div>
            <CurrentUser />
            <h1 className="text-2xl">{type} Chat with {name} | ID: {id}</h1>

            <div>
                {type === "private" ? (
                    <div>
                        {messages && messages.map((message) => (
                            <div key={message.id} className={`flex flex-col mb-2 ${String(message.sender.id ) === String(user?.id) ? "items-end" : "items-start"}`}>
                                <div className={`p-2 rounded-lg ${String(message.sender.id) === String(user?.id) ? "bg-blue-500 text-white" : "bg-[#4a4a4a]"}`}>
                                    <p>{message.content}</p>
                                    <small>{new Date(message.created_at).toLocaleString()}</small>
                                    {String(message.sender.id ) === String(user?.id) ? (<p>From: Me</p>) : (<p>From: {"recipient" in message && message.recipient.username}</p>)}
                                </div>
                            </div>
                        ))}
                    </div>
                ) : (
                    <div>
                        {messages && messages.map((message) => (
                            <div key={message.id} className={`flex flex-col mb-2 ${String(message.sender.id ) === String(user?.id) ? "items-end" : "items-start"}`}>
                                <div className={`p-2 rounded-lg ${String(message.sender.id) === String(user?.id) ? "bg-blue-500 text-white" : "bg-[#4a4a4a]"}`}>
                                    <p>{message.content}</p>
                                    <small>{new Date(message.created_at).toLocaleString()}</small>
                                    <p>From: {message.sender.username}</p>
                                </div>
                            </div>
                        ))}
                    </div>
                )}
                

                <input
                    type="text"
                    value={inputMessage}
                    onChange={(e) => setInputMessage(e.target.value)}
                    placeholder="Type your message..."
                    className="border border-gray-300 rounded-lg p-2"
                />
                <button onClick={sendMessage} className="bg-blue-500 text-white rounded-lg px-4 py-2 mt-2">Send</button>
            </div>
        </div>
    )
}

export default Chat
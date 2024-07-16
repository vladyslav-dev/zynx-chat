import { AuthContext } from '../providers/AuthProvider';
import { WebsocketContext } from '../providers/WebsocketProvider';
import React, { useState, useRef, useContext, useEffect } from 'react';
import { router } from "../providers/RouterProvider";
import { API_URL } from '../constants';
import { useLocation } from "react-router-dom";


export type Message = {
    content: string;
    client_id: string;
    username: string;
    room_id: string;
    type: 'recv' | 'self';
}

type UserState = { username: string }

const getClientsByRoomId = async (roomId: string) => {
    try {
        const response = await fetch(`${API_URL}/ws/getClients/${roomId}`, {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' },
        })
        const data = await response.json()

        return data
    } catch (err) {
        console.error(err)
    }

}

function Room() {
    const location = useLocation();
    const roomId = location.pathname.split('/').pop()
    console.log('roomId', roomId)

    // States
    const [users, setUsers] = useState<UserState[]>([])
    const [messages, setMessages] = useState<Message[]>([])

    // Refs
    const textareaRef = useRef<HTMLTextAreaElement>(null)

    // Context
    const { websocketConnection } = useContext(WebsocketContext)
    const { user } = useContext(AuthContext)

    useEffect(() => {
        if (websocketConnection === null) {
            router.navigate('/dashboard')
            return;
        }

        getClientsByRoomId(String(roomId)).then((data) => setUsers(data))
    }, [])

    useEffect(() => {
        if (websocketConnection === null) {
            router.navigate('/dashboard')
            return;
        }

        websocketConnection!.onmessage = (event) => {
            console.log('event', event)
            const message: Message = JSON.parse(event.data)

            if (message.content === 'A new user has joined the room') {
                setUsers(prevUsers => ([...prevUsers, { username: message.username }]))
            }

            if (message.content === 'User left the chat (Hub.Unregister case)') {
                setUsers(prevUsers => ([...prevUsers.filter(user => user.username !== message.username)]))
                setMessages(prevMessages => ([...prevMessages, message]))
                return
            }

            user?.username === message.username ? message.type = "self" : message.type = "recv"

            setMessages(prevMessages => ([...prevMessages, message]))

            websocketConnection!.onclose = (event) => {
                console.log('Connection closed', event)
            }
            websocketConnection!.onopen = (event) => {
                console.log('Connection opened', event)
            }
            websocketConnection!.onerror = (event) => {
                console.log('Connection error', event)
            }
        }
    }, [textareaRef, messages, websocketConnection, users])


    const sendMessage = () => {
        if (!textareaRef.current?.value.trim()) {
            alert('Please type a message')
            return;
        }

        if (websocketConnection === null) {
            router.navigate('/dashboard')
            return;
        }

        const message = textareaRef.current?.value
        websocketConnection?.send(message)

        textareaRef.current.value = ''
    }


    return (
        <div >
            <h2>Your current room is: {roomId}</h2>

            {messages.map((message, index) => {
                if (message.type === 'self') {
                    return (
                        <div
                            className='flex flex-col mt-2 w-full text-right justify-end bg-[#27292b]'
                            key={index}
                        >
                        <div className='text-sm'>{message.username} (me)</div>
                            <div>
                                <div className='bg-blue text-white px-4 py-1 rounded-md inline-block mt-1'>
                                    {message.content}
                                </div>
                            </div>
                        </div>
                    )
                } else {
                    return (
                        <div className='mt-2 bg-[#27292b]' key={index}>
                            <div className='text-sm'>{message.username}</div>
                            <div>
                                <div className='bg-grey text-dark-secondary px-4 py-1 rounded-md inline-block mt-1'>
                                    {message.content}
                                </div>
                            </div>
                        </div>
                    )
                }
            })}


            <textarea
                ref={textareaRef}
                placeholder='Type your message here'
                className='w-full h-10 p-2 rounded-md focus:outline-none'
                style={{ resize: 'none' }}
            />
            <button onClick={sendMessage}>Send</button>
        </div>
    );
}

export default Room;

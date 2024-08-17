import useGlobalStore from "../../shared/store/global"
import { Box, Flex, Heading, TabNav, Text } from "@radix-ui/themes"
import React, { useEffect, useRef, useState } from "react"
import { useLocation } from "react-router-dom"
import tempAvatar from "../../assets/temp-avatar.svg"
import useConversationStore, { GroupConversation, PrivateConversation } from "../../shared/store/conversation"
import useAuthStore from "../../shared/store/auth"
import { WEBSOCKET_URL } from "../../constants"
import CaretLeftIcon from "../../shared/components/icons/CaretLeftIcon"
import useScreenWidth from "../../shared/hooks/useScreenWidth"

const ConversationHeader = () => {
    const { conversationTitle, activeGroup, setActiveConversation } = useConversationStore()

    const onBackClick = () => {
        setActiveConversation(null)
    }

    return (
        <Flex className="items-center border-b border-zinc-300 w-full">
            <CaretLeftIcon onClick={onBackClick} className="inline-block w-[28px] h-[28px] cursor-pointer" />
            <Heading as={"h2"} className="p-4">{conversationTitle}</Heading>
        </Flex>
    )
}

const ConversationMessages = () => {
    const { user: me } = useAuthStore()
    const { messages, conversationMembers } = useConversationStore()

    const messagesEndRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (messages.length > 0) {
            messagesEndRef.current && messagesEndRef.current.scrollIntoView();
        }
        
    }, [messages]);

    return (
        <Flex direction="column" gap={"4"} className="p-4 overflow-auto">
            {messages.length ? messages.map((message) => (
                <Flex
                    key={message.id}
                    className="p-2"
                    style={{ justifyContent: me?.id === message.sender_id ? "flex-end" : "flex-start" }}
                >
                    {me?.id !== message.sender_id && (
                        <img
                            src={tempAvatar}
                            alt={"FILL THIS LATER"}
                            className="w-8 h-8 rounded-full mr-2"
                        />
                    )}
                    <Flex direction="column">
                        <Text size={"3"} weight={"medium"}>
                            {conversationMembers[message.sender_id]?.username || "Unknown"}
                        </Text>
                        <Text size={"3"} className="accent-text" style={{ background: "#F6F8FA", borderRadius: "8px", padding: "8px" }}>
                            {message.content}
                        </Text>
                        <Text size={"2"} className="opacity-50 mt-1">
                            {new Date(message.created_at).toLocaleString()}
                        </Text>
                    </Flex>
                    {me?.id === message.sender_id && (
                        <img
                            src={tempAvatar}
                            alt={"FILL THIS LATER"}
                            className="w-8 h-8 rounded-full ml-2"
                        />
                    )}
                </Flex>
            )) : <Text size={"4"} weight={"medium"} className="opacity-30">No messages yet</Text>}
            {/* This empty div ensures we can scroll to the bottom */}
            <div ref={messagesEndRef} />
        </Flex>
    )
}


export type MessagePayload = {
    type: string;
    sender_id: number;
    group_id?: number;
    recipient_id?: number;
    content: string;
}

export type PrivateMessagePayload = {
    type: "private";
    sender_id: number;
    recipient_id: number;
    content: string;
}

export type GroupMessagePayload = {
    type: "group";
    group_id: number;
    content: string;
}

const ConversationInput = () => {
    const { user: me } = useAuthStore()
    const { activeConversation, addMessage } = useConversationStore()
    const wsRef = useRef<WebSocket | null>(null)
    const [inputValue, setInputValue] = useState<string>("")

    useEffect(() => {
        if (!activeConversation) {
            console.error("No active conversation")
            return
        }

        if (wsRef.current) {
            wsRef.current.close()
        }

        const { type } = activeConversation;
        let URL = "";

        if (type === "group") {
            const { group_id } = activeConversation as GroupConversation;
            URL = `${WEBSOCKET_URL}/ws/message?type=group&group_id=${group_id}`;
        } else if (type === "private") {
            const { sender_id, recipient_id } = activeConversation as PrivateConversation;
            URL = `${WEBSOCKET_URL}/ws/message?type=private&sender_id=${sender_id}&recipient_id=${recipient_id}`;
        }

        wsRef.current = new WebSocket(URL)

        wsRef.current.onopen = () => {
            console.log('websocket connection opened')
        }

        wsRef.current.onmessage = (event) => {
            console.log('received message', JSON.parse(event.data))

            addMessage([JSON.parse(event.data)])
        }

        wsRef.current.onclose = () => {
            console.log('websocket connection closed')
        }

        wsRef.current.onerror = (error) => {
            console.error('websocket error', error)
        }

        return () => {
            console.log('closing websocket connection')
            wsRef.current?.close()
        }
    }, [activeConversation])

    const handleSendMessage = (event: React.FormEvent) => {
        event.preventDefault()

        if (!activeConversation) {
            console.error("No active conversation")
            return
        }

        let messagePayload

        if (activeConversation.type === "private") {
            const { sender_id, recipient_id } = activeConversation as PrivateConversation

            messagePayload = {
                type: "private",
                sender_id, recipient_id,
                content: inputValue
            } as PrivateMessagePayload
        }
        if (activeConversation?.type === "group") {
            const { group_id } = activeConversation as GroupConversation

            messagePayload = {
                type: "group",
                sender_id: me?.id,
                group_id: group_id,
                content: inputValue
            } as GroupMessagePayload
        }

        wsRef.current?.send(JSON.stringify(messagePayload))

        setInputValue("")
    }

    return (
        <Flex className="border-t border-zinc-300 w-full">
            <form onSubmit={handleSendMessage}>
                <Flex className="p-4">
                    <input 
                        type="text" 
                        value={inputValue} 
                        onChange={(event: React.ChangeEvent<HTMLInputElement>) => setInputValue(event.target.value)} 
                        name="message-input" 
                        className="w-full p-2 border border-zinc-300 rounded-md" 
                    />
                    <button type="submit" className="ml-2 bg-blue-500 text-white px-4 py-2 rounded-md">Send</button>
                </Flex>
            </form>
        </Flex>
    )
}

const ConversationDetails = () => {
    return (
        <Flex direction="column" className="relative border-x border-zinc-300 w-full h-full">
            <Box className="bg-white shadow-md z-10">
                <ConversationHeader />
            </Box>
            <Box className="flex-grow overflow-y-auto">
                <ConversationMessages />
            </Box>
            <Box className="bg-white shadow-md z-10">
                <ConversationInput />
            </Box>
        </Flex>
    );
}

const ConversationSettings = () => {


    return (
        <Flex direction="column" className="border-x border-zinc-300 w-[330px]">
            <Heading as={"h2"} className="p-4">Settings</Heading>
        </Flex>
    )
}

const ConversationView = () => {
    const { activeConversation } = useConversationStore()

    if (activeConversation) {
        return (
            <Flex className="overflow-scroll w-full">
                <ConversationDetails />
                {/* <ConversationSettings /> */}
            </Flex>
        )
    }

    return (
        <Flex className="h-full w-full" justify={"center"} align={"center"}>
            <Text size={"4"} weight={"medium"} className="opacity-30">Select a conversation</Text>
        </Flex>
    )
}

const PeopleList = () => {
    const { user: me } = useAuthStore()
    const { people } = useGlobalStore()
    const { setActiveConversation, activeConversation } = useConversationStore()

    const handleUserClick = (id: number) => {
        const ac = activeConversation as PrivateConversation
        if (ac && ac.type === "private" && ac.recipient_id === id) {
            /* Same user, do nothing */
            return
        }

        setActiveConversation({ type: "private", recipient_id: id, sender_id: me!.id })
    }

    return (
        <Flex direction="column">
            {people.map((user) => (
                <Flex
                    key={user.id}
                    className={`p-4 border-b border-zinc-300 cursor-pointer ${activeConversation?.type === "private" && (activeConversation as PrivateConversation).recipient_id === user.id ? "bg-blue-200" : ""}`}
                    onClick={() => handleUserClick(user.id)}
                >
                    <img src={tempAvatar} alt={user.username} className="w-8 h-8 rounded-full mr-2" />
                    <Text size={"3"} weight={"medium"}>{user.username}</Text>
                </Flex>
            ))}
        </Flex>
    )
}

const GroupList = () => {
    const { groups } = useGlobalStore()
    const { setActiveConversation, activeConversation } = useConversationStore()

    const onGroupClick = (id: number) => {
        const ac = activeConversation as GroupConversation
        if (ac && ac.type === "group" && ac.group_id === id) {
            /* Same group, do nothing */
            return
        }

        setActiveConversation({ type: "group", group_id: id })
    }

    return (
        <Flex direction="column">
            {groups.map((group) => (
                <Flex
                    key={group.id}
                    className={`p-4 border-b border-zinc-300 cursor-pointer ${activeConversation?.type === "group" && (activeConversation as GroupConversation).group_id === group.id ? "bg-blue-200" : ""}`}
                    onClick={() => onGroupClick(group.id)}
                >
                    <img src={"https://cdn1-production-images-kly.akamaized.net/GM1WEeVZycNwg_fGc_uSM0oDckA=/1200x1200/smart/filters:quality(75):strip_icc():format(webp)/kly-media-production/medias/915336/original/065075800_1435734425-Minions-1.jpg"} alt={group.name} className="w-8 h-8 rounded-full mr-2" />
                    <Text size={"3"} weight={"medium"}>{group.name}</Text>
                </Flex>
            ))}
        </Flex>
    )
}

const ConversationSidebar = () => {
    const location = useLocation()

    return (
        <Flex direction={"column"} className=" border-x border-zinc-300 h-full">
            <Heading as={"h2"} className="p-4">Chats</Heading>
            <TabNav.Root>
                <TabNav.Link href="#people" className="w-[160px]" active={location.hash === "#people" || location.hash === ""}>People</TabNav.Link>
                <TabNav.Link href="#groups" className="w-[160px]" active={location.hash === "#groups"}>Groups</TabNav.Link>
            </TabNav.Root>
            {location.hash === "#groups" ? <GroupList /> : <PeopleList />}
        </Flex>
    )
}

const Chat = () => {
    const { fetchPeople, fetchGroups } = useGlobalStore()
    const { activeConversation } = useConversationStore()
    const isMobile = useScreenWidth()

    useEffect(() => {
        fetchPeople()
        fetchGroups()
    }, [])

    return (
        <Flex className="h-dvh w-full">
            {isMobile ? (
                <>
                    {activeConversation ? <ConversationView /> : <ConversationSidebar />}
                </>            
            ) : (
                <>
                    <ConversationSidebar />
                    <ConversationView />
                </>
            )}
        </Flex>
    )
}

export default Chat
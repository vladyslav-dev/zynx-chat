import { createContext, useState } from "react";

interface IWebsocketContext {
    websocketConnection: WebSocket | null;
    setWebsocketConnection: (websocketConnection: WebSocket) => void;
}

const WebsocketContext = createContext<IWebsocketContext>({
    websocketConnection: null,
    setWebsocketConnection: (websocketConnection: WebSocket) => {}
})

const WebsocketProvider = ({ children }: { children: React.ReactNode }) => {
    const [websocketConnection, setWebsocketConnection] = useState<WebSocket | null>(null)

    const value = { websocketConnection, setWebsocketConnection }

    return (
        <WebsocketContext.Provider value={value}>
            {children}
        </WebsocketContext.Provider>
    )
}

export default WebsocketProvider;
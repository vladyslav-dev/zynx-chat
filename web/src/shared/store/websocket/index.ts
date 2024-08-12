import { create } from "zustand"

interface IWebsocketStore {
    websocketConnection: WebSocket | null;
    setWebsocketConnection: (ws: WebSocket) => void;
}

const useWebsocketStore = create<IWebsocketStore>((set, get) => ({
    websocketConnection: null,
    setWebsocketConnection: (ws) => {
        set({ websocketConnection: ws });
    },
}));

export default useWebsocketStore
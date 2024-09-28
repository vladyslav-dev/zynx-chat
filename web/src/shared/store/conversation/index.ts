import { create } from "zustand"
import { Message, MessageType } from '../../../shared/types/message';
import * as messageService from '../../../shared/api/message/message.service';
import * as groupService from '../../../shared/api/group/group.service';
import * as userService from '../../../shared/api/user/user.service';

export type PrivateConversation = {
    type: MessageType
    recipient_id: number
    sender_id: number
}

export type PrivateConversationPayload = {
    type: MessageType
    recipient_id: number
    sender_id: number
}

export type GroupConversation = {
    type: MessageType
    group_id: number
}

export type GroupConversationPayload = {
    type: MessageType
    group_id: number
}

export type Group = {
    id: number
    name: string
    createdAt: number
}

export type ConversationMember = {
    id: number
    username: string
    phone: string
}

export type ActiveConversation = PrivateConversation | GroupConversation | null;



interface IConversationStore {
    activeConversation: ActiveConversation;
    setActiveConversation: (conversation: ActiveConversation) => void;

    conversationTitle: string;
    setConversationTitle: (title: string) => void;

    activeGroup: Group | null;
    setActiveGroup: (group: Group) => void;
    // fetchActiveGroup: () => void;

    isGroupMember: boolean;
    setIsGroupMember: (isGroupMember: boolean) => void;

    /* One item for private conversation, multiple items for group conversation */
    conversationMembers: { [key: string]: ConversationMember };
    setConversationMembers: (members: any) => void;

    messages: Message[];
    setMessages: (messages: Message[]) => void;
    addMessage: (messages: Message[]) => void;
    fetchMessages: () => void;
}

const useConversationStore = create<IConversationStore>((set, get) => ({
    activeConversation: null,
    setActiveConversation: (conversation) => {
        set({ activeConversation: conversation, messages: [] });

        if (!conversation) {
            return;
        }

        get().fetchMessages();
    },

    conversationTitle: "",
    setConversationTitle: (title) => {
        set({ conversationTitle: title })
    },

    activeGroup: null,
    setActiveGroup: (group) => {
        set({ activeGroup: group })
    },
    // fetchActiveGroup: async () => {
    //     try {
    //         const { group_id } = get().activeConversation as GroupConversation;
    //         const response = await groupService.getGroupById(group_id);
    //         set({ activeGroup: response.data });
    //     } catch (error: any) {
    //         console.error(error);
    //         set({ activeGroup: null });
    //     }
    // },

    isGroupMember: false,
    setIsGroupMember: (isGroupMember) => {
        set({ isGroupMember: isGroupMember })
    },

    conversationMembers: {},
    setConversationMembers: (members) => {
        set({ conversationMembers: members })
    },

    messages: [],
    setMessages: (messages) => {
        set({ messages: messages })
    },
    addMessage: (messages) => {
        set({ messages: [...get().messages, ...messages] })
    },
    
    fetchMessages: () => {
        try {
            const type = get().activeConversation?.type

            if (!type) {
                console.error("No active conversation")
                
                return;
            }

            if (type === "private") {
                const { recipient_id, sender_id } = get().activeConversation as PrivateConversation

                const messagesPromise = messageService.getPrivateMessages({ type, recipient_id: recipient_id, sender_id: sender_id } as PrivateConversationPayload)
                const usersPromise = userService.getUsersByIds([recipient_id, sender_id]);

                Promise.all([messagesPromise, usersPromise])
                    .then(([messagesResponse, usersResponse]) => {
                        const members = Object.fromEntries(usersResponse.data.map((user: any) => [String(user.id), user]))
                        const conversationTitle = members[recipient_id].username
                        set({ conversationMembers: members, conversationTitle })

                        get().setMessages(messagesResponse.data)
                    })
                    .catch((error) => {
                        console.error(error)
                        set({ conversationMembers: {} })
                        get().setMessages([])
                    })
                    .finally(() => {

                        set({ activeGroup: null })
                    })
            } else if (type === "group") {
                const { group_id } = get().activeConversation as GroupConversation;

                const messagesPromise = messageService.getGroupMessages({ type, group_id: group_id } as GroupConversationPayload)
                const groupMembersPromise = userService.getUsersByGroupId(group_id)
                const groupPromise = groupService.getGroupById(group_id)

                Promise.all([messagesPromise, groupMembersPromise, groupPromise])
                    .then(([messagesResponse, groupMembersResponse, groupResponse]) => {
                        const members = Object.fromEntries(groupMembersResponse.data.map((user: any) => [user.id, user]))
                        const conversationTitle = groupResponse.data.name
                        set({ conversationMembers: members, activeGroup: groupResponse.data, conversationTitle });

                        get().setMessages(messagesResponse.data)
                    })
                    .catch((error) => {
                        console.error(error);
                        set({ conversationMembers: {}, activeGroup: null });
                        get().setMessages([])
                    });
            }
        } catch (error: any) {
            console.error(error);
            get().setMessages([])
        }
    }

}));

export default useConversationStore
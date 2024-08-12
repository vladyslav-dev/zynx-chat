import { GroupConversationPayload, PrivateConversationPayload } from "../../../shared/store/conversation"
import $api from "../index"

export const getPrivateMessages = async (payload: PrivateConversationPayload) => {
    return $api.post("/message/private", payload)
}

export const getGroupMessages = async (payload: GroupConversationPayload) => {
    return $api.post("/message/group", payload)
}
export type MessageType = 'private' | 'group'

export type PrivateMessage = {
    id: number;
    type: MessageType;
    sender_id: number;
    recipient_id: number;
    content: string;
    created_at: Date;
  }
  
  export type GroupMessage = {
    id: number;
    type: MessageType;
    sender_id: number;
    group_id: number;
    content: string;
    created_at: Date;
  }

export type Message = PrivateMessage | GroupMessage
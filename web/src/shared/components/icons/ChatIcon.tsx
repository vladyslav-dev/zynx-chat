import React, { ReactElement } from 'react';
import chatBuble from '../../../assets/icons/chat-bubble.svg';

const ChatIcon = (props: React.ImgHTMLAttributes<HTMLImageElement>): ReactElement => <img src={chatBuble} alt='chat-buble' {...props} />

export default ChatIcon
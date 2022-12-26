import React from 'react';
import {AiOutlineEdit, AiOutlineDelete} from "react-icons/ai";
import {MsgData, editMessage, deleteMessage} from "./Message";

export type MessageButtonProps = {
    msg: MsgData
    data: MsgData[]
    setData: React.Dispatch<React.SetStateAction<MsgData[]>>
}


export const MessageButtons: React.FC<MessageButtonProps> = (props: MessageButtonProps) => {

    return (<div style={{display: 'flex', justifyContent: 'flex-end'}}>
        <button onClick={() => {
            editMessage(props.msg, props.data, props.setData)
        }}><AiOutlineEdit/></button>
        <button onClick={() => {
            deleteMessage(props.msg.id, props.data, props.setData)
        }}><AiOutlineDelete/></button>
    </div>);
}